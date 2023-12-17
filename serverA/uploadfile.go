package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func (sc *ConfigFile) splitAndUpload(r *http.Request) error {
	//Получаем файл через multipart/form и парсим и анализируем форму.
	err := r.ParseMultipartForm(10000000)
	if err != nil {
		return fmt.Errorf("can not parse form:, %v", err)
	}
	fileData, fileHead, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("no file field in request: %v", err)
	}
	defer fileData.Close()

	if len(fileHead.Filename) <= 0 {
		return errors.New("incorrect length filename")
	}

	//Вычисляем первый сервер с помощью генератора случайных чисел из диапазона от 0 до количества серверов
	firstSrv := rand.Intn(sc.CountSrv)

	var wg sync.WaitGroup
	errChan := make(chan bool, 6)
	offset := fileHead.Size / 6 //вычисляем размер одной части файла

	var srvNum int

	for i := 0; i <= 5; i++ {
		//режем файл в соответствии с вычисленным размером части
		buf := make([]byte, offset, offset+fileHead.Size%6)
		//в последнюю часть складываем все оставшиеся байты.
		if i == 5 && fileHead.Size%6 > 0 {
			buf = buf[:(offset + fileHead.Size%6)]
		}
		n, err := fileData.ReadAt(buf, offset*int64(i))
		//вычисляем сервер для загрузки
		srvNum = firstSrv + i
		if srvNum > sc.CountSrv-1 {
			srvNum = srvNum - sc.CountSrv
		}
		//Отдаем данные в горутину для загрузки
		wg.Add(1)
		go sc.uploadToB(&wg, errChan, buf[:n], i, fileHead.Filename, srvNum)
		if err == io.EOF {
			break
		}
	}
	wg.Wait()
	close(errChan)
	// проверяем, что все части получены без ошибки
	for errVal := range errChan {
		if errVal == false {
			return errors.New("upload to serverB failed")
		}
	}
	// создаем каталог для метафайла
	err = os.MkdirAll(filepath.Join(sc.BaseDir, string(fileHead.Filename[0])), os.ModePerm)
	if err != nil {
		return fmt.Errorf("can not create directory: %v", err)
	}
	// считаем контрольную сумму полученного файла
	h := sha256.New()
	if _, err := io.Copy(h, fileData); err != nil {
		return fmt.Errorf("hash calc error: %v", err)
	}
	// записываем в метафайл контрольную сумму и сервер, на который загружена первая часть файла
	metaFile, _ := json.Marshal(&fileInfo{
		CrcSum:      hex.EncodeToString(h.Sum(nil)),
		FirstServer: firstSrv,
	})

	err = os.WriteFile(filepath.Join(sc.BaseDir, string(fileHead.Filename[0]), fileHead.Filename), metaFile, 0644)
	if err != nil {
		return fmt.Errorf("can not write meta file: %v", err)
	}

	return nil

}

func (sc *ConfigFile) uploadToB(wg *sync.WaitGroup, errChan chan<- bool, uploadChunk []byte, idxChunk int, filename string, srvnum int) {
	defer wg.Done()
	buf := bytes.NewReader(uploadChunk)
	resp, err := http.Post(fmt.Sprintf("http://%s/%c/%s.%d", sc.StorageServers[srvnum], filename[0], filename, idxChunk), "application/octet-stream", buf)
	if err != nil {
		errChan <- false
		return
	}
	if resp.StatusCode != 200 {
		errChan <- false
		return
	}
	errChan <- true
	return
}
