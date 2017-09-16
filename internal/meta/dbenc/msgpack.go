package dbenc

import (
	"encoding/hex"
	"time"

	"github.com/tinylib/msgp/msgp"

	"github.com/ophymx/s3d/internal/meta"
)

type msgpEncoding struct {
}

var MsgPack = msgpEncoding{}

type bucketField uint8

const (
	bucketCreation bucketField = 1
)

func (e msgpEncoding) EncodeBucket(data meta.BucketData) (b []byte, err error) {
	b = msgp.AppendMapHeader(nil, 1)
	b = e.appendBucketField(b, bucketCreation)
	b = e.appendTime(b, data.CreationDate)
	return
}

func (e msgpEncoding) DecodeBucket(b []byte) (data meta.BucketData, err error) {
	sz, b, err := msgp.ReadMapHeaderBytes(b)
	if err != nil {
		return
	}
	for i := uint32(0); i < sz; i++ {
		var field bucketField
		if field, b, err = e.readBucketField(b); err != nil {
			return
		}
		switch field {
		case bucketCreation:
			data.CreationDate, b, err = e.readTime(b)
		}
		if err != nil {
			return
		}
	}
	return
}

type objectField uint8

const (
	objectContentMD5   objectField = 1
	objectSize                     = 2
	objectCacheControl             = 3
	objectLastModified             = 4
	objectContentType              = 5
	objectVersionID                = 6
	objectUserDefined              = 7
)

func (e msgpEncoding) EncodeObject(data meta.ObjectData) (b []byte, err error) {
	b = msgp.AppendMapHeader(nil, 7)

	b = e.appendObjectField(b, objectContentMD5)
	md5, err := hex.DecodeString(data.ContentMD5)
	if err != nil {
		return
	}
	b = msgp.AppendBytes(b, md5)

	b = e.appendObjectField(b, objectSize)
	b = msgp.AppendInt64(b, data.Size)

	b = e.appendObjectField(b, objectCacheControl)
	b = msgp.AppendString(b, data.CacheControl)

	b = e.appendObjectField(b, objectLastModified)
	b = e.appendTime(b, data.LastModified)

	b = e.appendObjectField(b, objectContentType)
	b = msgp.AppendString(b, data.ContentType)

	b = e.appendObjectField(b, objectVersionID)
	b = msgp.AppendString(b, data.VersionID)

	b = e.appendObjectField(b, objectUserDefined)
	b = msgp.AppendMapHeader(b, uint32(len(data.UserDefined)))
	if data.UserDefined != nil {
		for k, v := range data.UserDefined {
			b = msgp.AppendString(b, k)
			b = msgp.AppendString(b, v)
		}
	}

	return
}

func (e msgpEncoding) DecodeObject(b []byte) (data meta.ObjectData, err error) {
	var sz, i uint32
	if sz, b, err = msgp.ReadMapHeaderBytes(b); err != nil {
		return
	}
	for i = 0; i < sz; i++ {
		var field objectField
		if field, b, err = e.readObjectField(b); err != nil {
			return
		}
		switch field {
		case objectContentMD5:
			var md5 []byte
			if md5, b, err = msgp.ReadBytesBytes(b, nil); err == nil {
				data.ContentMD5 = hex.EncodeToString(md5)
			}
		case objectSize:
			data.Size, b, err = msgp.ReadInt64Bytes(b)
		case objectCacheControl:
			data.CacheControl, b, err = msgp.ReadStringBytes(b)
		case objectLastModified:
			data.LastModified, b, err = e.readTime(b)
		case objectContentType:
			data.ContentType, b, err = msgp.ReadStringBytes(b)
		case objectVersionID:
			data.VersionID, b, err = msgp.ReadStringBytes(b)
		case objectUserDefined:
			data.UserDefined, b, err = e.readMapStrStr(b)
		}
		if err != nil {
			return
		}
	}
	return
}

func (e msgpEncoding) readMapStrStr(in []byte) (m map[string]string, b []byte, err error) {
	var sz, i uint32
	if sz, b, err = msgp.ReadMapHeaderBytes(in); err != nil {
		return
	}
	m = make(map[string]string, int(sz))
	for i = 0; i < sz; i++ {
		var k, v string
		if k, b, err = msgp.ReadStringBytes(b); err != nil {
			return
		}
		if v, b, err = msgp.ReadStringBytes(b); err != nil {
			return
		}
		m[k] = v
	}
	return
}

func (e msgpEncoding) appendBucketField(b []byte, field bucketField) []byte {
	return msgp.AppendUint8(b, uint8(field))
}

func (e msgpEncoding) readBucketField(b []byte) (bucketField, []byte, error) {
	i, b, err := msgp.ReadUint8Bytes(b)
	if err != nil {
		return 0, b, err
	}
	return bucketField(i), b, err
}

func (e msgpEncoding) appendObjectField(b []byte, field objectField) []byte {
	return msgp.AppendUint8(b, uint8(field))
}

func (e msgpEncoding) readObjectField(b []byte) (objectField, []byte, error) {
	i, b, err := msgp.ReadUint8Bytes(b)
	if err != nil {
		return 0, b, err
	}
	return objectField(i), b, err
}

func (e msgpEncoding) appendTime(b []byte, t time.Time) []byte {
	b, _ = msgp.AppendExtension(b, newTime(t))
	return b
}

func (e msgpEncoding) readTime(b []byte) (time.Time, []byte, error) {
	ext := &timeExtension{}
	b, err := msgp.ReadExtensionBytes(b, ext)
	return ext.UTC(), b, err
}
