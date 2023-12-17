package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {

	//Загрузка файла

	multipartBuf := bytes.NewBuffer(nil)

	multipartWriter := multipart.NewWriter(multipartBuf)
	fw, err := multipartWriter.CreateFormFile("file", "ava_tux.jpg")
	if err != nil {
		log.Println(err)
	}

	buf, err := os.ReadFile("client/ava_tux.jpg")
	if err != nil {
		log.Println(err)
	}
	data := bytes.NewReader(buf)
	_, err = io.Copy(fw, data)
	if err != nil {
		log.Println(err)
	}

	err = multipartWriter.Close()
	if err != nil {
		log.Println(err)
	}

	sendFileRequest, err := http.NewRequest(http.MethodPost, "http://localhost:8080", multipartBuf)
	sendFileRequest.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	if err != nil {
		log.Println(err)
	}

	httpClient := http.Client{}
	sendFileResponse, err := httpClient.Do(sendFileRequest)
	if err != nil {
		log.Println(err)
	}

	response, err := io.ReadAll(sendFileResponse.Body)
	if err != nil {
		log.Println(err)
	}

	log.Println(sendFileResponse.StatusCode, string(response))

	//Получение файла
	getFileResponce, err := http.Get("http://localhost:8080/ava_tux.jpg")
	if err != nil {
		log.Println(err)
	}
	log.Println(getFileResponce.StatusCode)
	respData, err := io.ReadAll(getFileResponce.Body)
	err = os.WriteFile(os.TempDir()+"/ava_tux.jpg", respData, 0644)
	if err != nil {
		log.Println(err)
	}

}
