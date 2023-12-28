package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ConfigFile struct {
	Listen  string `json:"listen"`
	BaseDir string `json:"base_dir"`
}

func parseConfig(configPath string) (*ConfigFile, error) {
	fileBody, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var sc ConfigFile
	err = json.Unmarshal(fileBody, &sc)
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

func (sc *ConfigFile) serverBHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		//Запись файла, данные берем из тела запроса POST
		if r.ContentLength == 0 {
			log.Println("Empty post request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dir, _ := filepath.Split(r.URL.Path)
		err := os.MkdirAll(filepath.Join(sc.BaseDir, dir), os.ModePerm)
		if err != nil {
			log.Println("Can not create directory")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		file, err := os.OpenFile(filepath.Join(sc.BaseDir, r.URL.Path), os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Println("Can not create new file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, r.Body)
		if err != nil {
			log.Println("Can not write file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		//Отдача файла
		http.ServeFile(w, r, filepath.Join(sc.BaseDir, r.URL.Path))
	default:
		http.NotFound(w, r)
	}
}

func main() {
	confPath := flag.String("config", "serverB/config.json", "")
	flag.Parse()
	sConfig, err := parseConfig(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", sConfig.serverBHandler)
	err = http.ListenAndServe(sConfig.Listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}
