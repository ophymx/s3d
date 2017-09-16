package server

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ophymx/s3d/internal/auth"
	"github.com/ophymx/s3d/internal/blob"
	"github.com/ophymx/s3d/internal/clock"
	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/ops"
	"github.com/ophymx/s3d/internal/s3"
)

const (
	MethodGET    = "GET"
	MethodHEAD   = "HEAD"
	MethodPUT    = "PUT"
	MethodPOST   = "POST"
	MethodDELETE = "DELETE"
)

type Service interface {
	Serve(req s3.Request) (res s3.Response)
}

type S3Handler struct {
	objectService  Service
	bucketService  Service
	serviceService Service
	bucketParser   s3.BucketParser
	credentials    map[string]s3.Credential
	config         s3.Config
	requestID      uint64
}

func NewHandler(
	db meta.DB,
	store blob.Store,
	bucketParser s3.BucketParser,
	credentials map[string]s3.Credential,
	config s3.Config,
) http.Handler {

	return &S3Handler{
		bucketService:  NewBucketService(ops.NewBucket(db, store, clock.Real)),
		objectService:  NewObjectService(ops.NewObject(db, store, clock.Real)),
		serviceService: NewServiceService(ops.NewService(db)),
		bucketParser:   bucketParser,
		credentials:    credentials,
		config:         config,
		requestID:      uint64(time.Now().Unix()),
	}
}

func (h *S3Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	start := time.Now()
	requestID := h.getRequestID()
	resource := h.getResource(req.Host, req.URL.Path)
	response := h.serve(start, requestID, resource, req)
	if response == nil {
		response = s3.InternalErrorf("no response")
	}

	writer.Header().Add(s3.AmzRequestID, requestID)
	writer.Header().Add(s3.AmzHostID, h.config.HostID)
	writer.Header().Add("Server", "s3d")
	err := response.Send(writer)
	log.Printf(
		"[%s] (%s) %s %s %s %d %vµs",
		requestID,
		resource.Bucket(),
		req.RemoteAddr,
		req.Method,
		req.URL,
		response.HTTPStatus(),
		int64(time.Since(start)/time.Microsecond),
	)
	if err != nil {
		log.Printf("[%s] Error Sending: %s", requestID, err.Error())
	}
}

func (h *S3Handler) serve(
	start time.Time,
	requestID string,
	resource s3.Resource,
	req *http.Request,
) (response s3.Response) {
	values, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return s3.InvalidRequest("invalid query")
	}

	authorization, err := auth.GetAuth(resource, req.Header, req.URL.Query())
	if err != nil {
		return s3.AuthError(err)
	}

	var cred s3.Credential
	if authorization != nil {
		var found bool
		accessKeyID := authorization.GetAccessKeyID()
		if cred, found = h.credentials[accessKeyID]; !found {
			return s3.InvalidAccessKeyID(accessKeyID)
		}

		if err = authorization.Verify(cred.SecretKey, req); err != nil {
			return s3.AuthError(err)
		}
	}

	return h.route(resource).Serve(s3.Request{
		ID:         requestID,
		Method:     req.Method,
		Auth:       authorization,
		Credential: cred,
		Resource:   resource,
		Query:      values,
		Time:       start,
		RawReq:     req,
	})
}

func (h *S3Handler) getRequestID() string {
	return strings.ToUpper(strconv.FormatUint(atomic.AddUint64(&h.requestID, 1), 16))
}

func (h *S3Handler) route(resource s3.Resource) Service {
	switch {
	case resource.Key() != "":
		return h.objectService
	case resource.Bucket() != "":
		return h.bucketService
	default:
		return h.serviceService
	}
}

func (h *S3Handler) getResource(host, path string) (resource s3.Resource) {
	path = strings.TrimLeft(path, "/")
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[0:idx]
	}

	if bucket := h.bucketParser.Parse(host); bucket != "" {
		return s3.NewResource(bucket, path)
	}
	return s3.ParseResource(path)
}

type ObjectService struct {
	ops.ObjectOperations
}

func NewObjectService(ops ops.ObjectOperations) ObjectService {
	return ObjectService{ops}
}

func (srv ObjectService) Serve(req s3.Request) s3.Response {
	switch req.Method {
	case MethodGET:
		return srv.Get(req.Resource)
	case MethodHEAD:
		return srv.Head(req.Resource)
	case MethodDELETE:
		return srv.Delete(req.Resource)
	case MethodPOST:
		return s3.MethodNotAllowed(req.Method + " method not allowed on bucket")
	case MethodPUT:
		if copySrc := req.RawReq.Header.Get(s3.AmzCopySource); copySrc != "" {
			return srv.Copy(s3.ParseResource(copySrc), req.Resource)
		}
		return srv.Put(req.Resource, req.RawReq.Header.Get(s3.HdrContentType), req.RawReq.Body)
	default:
		return s3.MethodNotAllowed(req.Method + " method not allowed on bucket")
	}
}

type BucketService struct {
	ops.BucketOperations
}

func NewBucketService(ops ops.BucketOperations) BucketService {
	return BucketService{ops}
}

func (srv BucketService) Serve(req s3.Request) s3.Response {
	switch req.Method {
	case MethodGET:
		return srv.ListBucket(req.Resource.Bucket(), req.Query)
	case MethodHEAD:
		return srv.ListBucket(req.Resource.Bucket(), req.Query)
	case MethodPUT:
		return srv.Create(req.Resource.Bucket())
	case MethodDELETE:
		return srv.Delete(req.Resource.Bucket())
	default:
		return s3.MethodNotAllowed(req.Method + " method not allowed on bucket")
	}
}

type ServiceService struct {
	ops.ServiceOperations
}

func NewServiceService(ops ops.ServiceOperations) (service ServiceService) {
	return ServiceService{ops}
}

func (srv ServiceService) Serve(req s3.Request) s3.Response {
	switch req.Method {
	case MethodGET:
		return srv.ListBuckets()
	case MethodHEAD:
		return srv.ListBuckets()
	default:
		return s3.MethodNotAllowed(req.Method + " not allowed on service")
	}
}
