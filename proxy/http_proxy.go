package http_cache

import (
	"bufio"
	"context"
	"fmt"
	"cloud.google.com/go/storage"
	"log"
	"net"
	"net/http"
)

type StorageProxy struct {
	bucketHandler *storage.BucketHandle
	defaultPrefix string
}

func NewStorageProxy(bucketHandler *storage.BucketHandle, defaultPrefix string) *StorageProxy {
	return &StorageProxy{
		bucketHandler: bucketHandler,
		defaultPrefix: defaultPrefix,
	}
}

func (proxy StorageProxy) objectName(name string) string {
	return proxy.defaultPrefix + name
}

func (proxy StorageProxy) Serve(port int64) error {
	http.HandleFunc("/", proxy.handler)

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))

	if err == nil {
		address := listener.Addr().String()
		listener.Close()
		log.Printf("Starting http cache server %s\n", address)
		return http.ListenAndServe(address, nil)
	}
	return err
}

func (proxy StorageProxy) handler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	if key[0] == '/' {
		key = key[1:]
	}
	if r.Method == "GET" {
		proxy.downloadBlob(w, key)
	} else if r.Method == "HEAD" {
		proxy.checkBlobExists(w, key)
	} else if r.Method == "POST" {
		proxy.uploadBlob(w, r, key)
	} else if r.Method == "PUT" {
		proxy.uploadBlob(w, r, key)
	}
}

func (proxy StorageProxy) downloadBlob(w http.ResponseWriter, name string) {
	object := proxy.bucketHandler.Object(proxy.objectName(name))
	if object == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	reader, err := object.NewReader(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer reader.Close()
	bufferedReader := bufio.NewReader(reader)
	_, err = bufferedReader.WriteTo(w)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (proxy StorageProxy) checkBlobExists(w http.ResponseWriter, name string) {
	object := proxy.bucketHandler.Object(proxy.objectName(name))
	if object == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// lookup attributes to see if the object exists
	attrs, err := object.Attrs(context.Background())
	if err != nil || attrs == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (proxy StorageProxy) uploadBlob(w http.ResponseWriter, r *http.Request, name string) {
	object := proxy.bucketHandler.Object(proxy.objectName(name))

	writer := object.NewWriter(context.Background())
	defer writer.Close()

	_, err := bufio.NewWriter(writer).ReadFrom(bufio.NewReader(r.Body))
	if err != nil {
		errorMsg := fmt.Sprintf("Failed read cache body! %s", err)
		w.Write([]byte(errorMsg))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
