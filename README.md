S3 Daemon (For Testing)
=======================

`s3d` is a partial implementation of an [S3](https://aws.amazon.com/s3/) server trying to follow as closely
to the [REST API](http://docs.aws.amazon.com/AmazonS3/latest/API/Welcome.html) as possible.
The intent is to use it in local testing, as an alternative to making actual calls to AWS.

Warning!!
---------
This implementation of an S3 server is for testing purposes only.
__Do not use this with production data!__
Reads and writes don't actually require authentication,
no quotas of any kind are enforced, and all input validation
is focused on mimicing responses from AWS S3 and not focused on actual security.

Database and file storage layout are not stable and will probably change in the future.
The `s3d` data directory should probably be deleted between test runs.

Installing and Running
-------

```
go get github.com/ophymx/s3d
s3d -d DATA_DIR -p PORT -a TEST_ACCESS_ID -s TEST_SECRET_KEY
```
*NOTE: Don't use actual AWS credentials with this.*

What's implemented so far?
--------------------------
- List Buckets
- Create and Delete Bucket
- PUT, GET, HEAD, DELETE Object and Copy via PUT
- List Objects in Bucket
- Authentication with V2 and V4 signatures
  - Using either HTTP Header or URL Query authentication
  - Only single chunk so far
  - Authentication is only validated if present in a request. It is not enforced.

See [Issues](https://github.com/ophymx/s3d/issues?utf8=%E2%9C%93&q=is%3Aissue%20label%3Aenhancement)
to track additional features.

Alternatives
------------

`s3d` is pretty young and has not seen much use yet.
If you're looking for an alternative, consider [fakes3](https://github.com/jubos/fake-s3).
