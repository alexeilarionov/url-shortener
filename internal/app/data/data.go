package data

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

var (
	dataMap        = make(map[string]string)
	mutex          = sync.RWMutex{}
	dataFilePath   = "/Users/lari4/go/src/url-shortener/storage/data.txt"
	ErrKeyNotFound = errors.New("key not found")
)

func InitData() {
	err := loadDataFromFile()
	if err != nil {
		log.Fatalf("Failed to load data from file: %v", err)
	}
}

func GetAllData() (map[string]string, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	response := make(map[string]string)
	for k, v := range dataMap {
		response[k] = v
	}
	return response, nil
}

func GetDataByKey(key string) (string, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	value, found := dataMap[key]
	if !found {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func AddOrUpdateData(k string, v string) error {
	mutex.Lock()
	defer mutex.Unlock()

	dataMap[k] = v

	err := saveDataToFile()
	if err != nil {
		return err
	}

	return nil
}

func loadDataFromFile() error {
	dataFile, err := os.Open(dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer dataFile.Close()

	err = json.NewDecoder(dataFile).Decode(&dataMap)
	if err != nil {
		return err
	}

	return nil
}

func saveDataToFile() error {
	dataFile, err := os.OpenFile(dataFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	err = json.NewEncoder(dataFile).Encode(dataMap)
	if err != nil {
		return err
	}

	return nil
}
