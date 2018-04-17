package main

import (
	"cloud.google.com/go/storage"
	"context"
	"flag"
	"github.com/cirruslabs/google-storage-proxy/proxy"
	"log"
)

func main() {
	var port int64
	flag.Int64Var(&port, "port", 8080, "Port to serve")
	var bucketName string
	flag.StringVar(&bucketName, "bucket", "", "Google Storage Bucket Name")
	var defaultPrefix string
	flag.StringVar(&defaultPrefix, "prefix", "", "Default prefix for all object names. For example, use --prefix=foo/.")
	flag.Parse()

	if bucketName == "" {
		log.Fatal("Please specify Google Cloud Storage Bucket")
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create a storage client: %s", err)
	}

	bucketHandler := client.Bucket(bucketName)
	storageProxy := http_cache.NewStorageProxy(bucketHandler, defaultPrefix)

	err = storageProxy.Serve(port)
	if err != nil {
		log.Fatalf("Failed to start proxy: %s", err)
	}
}
