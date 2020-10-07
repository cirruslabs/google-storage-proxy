package main

import (
	"cloud.google.com/go/storage"
	"context"
	"flag"
	"github.com/cirruslabs/google-storage-proxy/proxy"
	"log"
)

func main() {
	var address string
	flag.StringVar(&address, "address", "127.0.0.1", "Address to listen on")
	var port int64
	flag.Int64Var(&port, "port", 8080, "Port to serve")
	var bucketName string
	flag.StringVar(&bucketName, "bucket", "", "Google Storage Bucket Name")
	var defaultPrefix string
	flag.StringVar(&defaultPrefix, "prefix", "", "Optional prefix for all objects. For example, use --prefix=foo/ to work under foo directory in a bucket.")
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

	err = storageProxy.Serve(address, port)
	if err != nil {
		log.Fatalf("Failed to start proxy: %s", err)
	}
}
