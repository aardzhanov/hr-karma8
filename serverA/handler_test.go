package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetHandler(t *testing.T) {
	sConfig, err := parseConfig("config.json")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/ava_tux.jpg", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sConfig.userApiHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if rr.Body.Len() != 19577 {
		t.Errorf("handler returned unexpected body: got %d want 19577", rr.Body.Len())
	}
}

func TestPostHandler(t *testing.T) {
	sConfig, err := parseConfig("config.json")
	if err != nil {
		t.Fatal(err)
	}

	multipartBuf := bytes.NewBuffer(nil)

	multipartWriter := multipart.NewWriter(multipartBuf)
	fw, err := multipartWriter.CreateFormFile("file", "ava_tux.jpg")
	if err != nil {
		t.Fatal(err)

	}

	buf, err := os.ReadFile("../client/ava_tux.jpg")
	if err != nil {
		t.Fatal(err)

	}
	data := bytes.NewReader(buf)
	_, err = io.Copy(fw, data)
	if err != nil {
		t.Fatal(err)
	}

	err = multipartWriter.Close()
	if err != nil {
		t.Fatal(err)
	}

	sendFileRequest, err := http.NewRequest(http.MethodPost, "/", multipartBuf)
	sendFileRequest.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	if err != nil {
		t.Fatal(err)

	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sConfig.userApiHandler)

	handler.ServeHTTP(rr, sendFileRequest)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestPutHandler(t *testing.T) {
	sConfig, err := parseConfig("config.json")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "/ava_tux.jpg", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sConfig.userApiHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}

}
