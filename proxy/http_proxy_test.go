package http_cache

import (
	"testing"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"net/http/httptest"
	"net/http"
	"strings"
)

const TestBucketName = "some-bucket"

func Test_Blob_Exists(t *testing.T) {
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			BucketName: TestBucketName,
			Name:       "some/object/file",
		},
	})
	defer server.Stop()
	client := server.Client()
	storageProxy := NewStorageProxy(client.Bucket(TestBucketName), "")

	response := httptest.NewRecorder()
	storageProxy.checkBlobExists(response, "some/object/file")

	if response.Code == http.StatusOK {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d'", response.Code)
	}
}

func Test_Default_Prefix(t *testing.T) {
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			BucketName: TestBucketName,
			Name:       "some/object/file",
		},
	})
	defer server.Stop()
	client := server.Client()
	storageProxy := NewStorageProxy(client.Bucket(TestBucketName), "some/object/")

	response := httptest.NewRecorder()
	storageProxy.checkBlobExists(response, "file")

	if response.Code == http.StatusOK {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d'", response.Code)
	}
}

func Test_Blob_Download(t *testing.T) {
	expectedBlobContent := "my content"
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			BucketName: TestBucketName,
			Name:       "some/file",
			Content:    []byte(expectedBlobContent),
		},
	})
	defer server.Stop()
	client := server.Client()
	storageProxy := NewStorageProxy(client.Bucket(TestBucketName), "")

	response := httptest.NewRecorder()
	storageProxy.downloadBlob(response, "some/file")

	if response.Code == http.StatusOK {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d'", response.Code)
	}

	downloadedBlobContent := response.Body.String()
	if downloadedBlobContent == expectedBlobContent {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong content: '%s'", downloadedBlobContent)
	}
}

func Test_Blob_Upload(t *testing.T) {
	expectedBlobContent := "my content"
	server := fakestorage.NewServer([]fakestorage.Object{})
	defer server.Stop()
	client := server.Client()
	storageProxy := NewStorageProxy(client.Bucket(TestBucketName), "")

	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/test-file", strings.NewReader(expectedBlobContent))
	storageProxy.uploadBlob(response, request,"test-file")

	if response.Code == http.StatusCreated {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d'", response.Code)
	}
}