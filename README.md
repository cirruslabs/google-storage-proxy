[![Build Status](https://api.cirrus-ci.com/github/cirruslabs/google-storage-proxy.svg)](https://cirrus-ci.com/github/cirruslabs/google-storage-proxy)

HTTP proxy with REST API to interact with Google Cloud Storage Buckets

Simply allows using `HEAD`, `GET` or `PUT` requests to check blob's availability, as well as downloading or uploading
blobs to a specified GCS bucket.

# Arguments

* `port` - optional port to run the proxy on. By default, `8080` is used.
* `bucket` - GCS bucket name to store artifacts in. You can configure [Lifecycle Management](https://cloud.google.com/storage/docs/lifecycle)
   for this bucket separately using `gcloud` or UI.
* `prefix` - optional prefix for all objects. For example, use `--prefix=foo/` to work under `foo` directory in `bucket`.
