package main

import (
	"encoding/json"
	"os"
)

// Функция считывания файла конфигурации в формате json в структуру
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
	sc.CountSrv = len(sc.StorageServers)
	return &sc, nil
}
