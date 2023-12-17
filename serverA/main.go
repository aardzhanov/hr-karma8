package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
)

func (sc *ConfigFile) userApiHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		//Загрузка файла при обращении методом POST
		err := sc.splitAndUpload(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Println(err)
			}
			return
		}
	case http.MethodGet:
		//Получение файла при обращении методом GET
		file, err := sc.combineAndDownload(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Println(err)
			}
			return
		}
		_, err = w.Write(file.Bytes())
		if err != nil {
			log.Println(err)
		}
	default:
		//В случе остальных методов выдаем MethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	//Считываем переданный флаг с путем к файлу конфигурации
	confPath := flag.String("config", "serverA/config.json", "")
	flag.Parse()
	//Читаем файл конфигурации
	sConfig, err := parseConfig(*confPath)
	if err != nil {
		log.Fatal(err)
	}
	//Если в конфигурации серверов меньше 6, то логируем ошибку и выходим
	if len(sConfig.StorageServers) < 6 {
		log.Fatal(errors.New("servers count less than 6"))
	}

	http.HandleFunc("/", sConfig.userApiHandler)
	err = http.ListenAndServe(sConfig.Listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}
