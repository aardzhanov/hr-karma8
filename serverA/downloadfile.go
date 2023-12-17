package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Функция скачивания частей файла, их склеивания
func (sc *ConfigFile) combineAndDownload(r *http.Request) (*bytes.Buffer, error) {
	//вырезаем все пути переданные пользователем и берем только имя файла
	_, fileName := filepath.Split(r.URL.Path)

	//читаем метафайл, хранящий контрольную сумм и индекс первого сервера
	fileBody, err := os.ReadFile(filepath.Join(sc.BaseDir, string(fileName[0]), fileName))
	if err != nil {
		return nil, fmt.Errorf("can not read meta file: %v", err)
	}
	var metaData fileInfo

	err = json.Unmarshal(fileBody, &metaData)
	if err != nil {
		return nil, fmt.Errorf("can not unmarshal meta file: %v", err)
	}

	var srvNum int

	//делаем запрос к серверам Б и получаем все части файла
	var bufFile bytes.Buffer
	for i := 0; i <= 5; i++ {
		//считаем сервер на котором хранится часть файла (сервера идут по порядку циклично)
		srvNum = metaData.FirstServer + i
		if srvNum > sc.CountSrv-1 {
			srvNum = srvNum - sc.CountSrv
		}
		// делаем запрос к серверу Б и анализируем ответ
		resp, err := http.Get(fmt.Sprintf("http://%s/%c/%s.%d", sc.StorageServers[srvNum], fileName[0], fileName, i))
		if err != nil {
			return nil, fmt.Errorf("can not get chunk %d: %v", i, err)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status for chunk %d not 200: %d (%s)", i, resp.StatusCode, resp.Status)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("can not read body of chunk %d: %v", i, err)
		}
		//дописываем прочитанную часть в буфер
		_, err = bufFile.Write(body)
		if err != nil {
			return nil, fmt.Errorf("can not write chunk %d to buffer: %v", i, err)
		}

	}
	// расчет контрольной суммы в буфера
	h := sha256.New()
	_, err = h.Write(bufFile.Bytes())
	if err != nil {
		return nil, fmt.Errorf("hash calc error: %v", err)

	}
	//  сравнение контрольных сумм в буфере и метафайле
	if hex.EncodeToString(h.Sum(nil)) != metaData.CrcSum {
		return nil, errors.New("crc error")
	}

	return &bufFile, nil
}
