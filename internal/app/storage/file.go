package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type FileStorage struct {
	inMemoryData *InMemoryStorage
	filePath     string
}

func NewFileStorage(filename string) *FileStorage {
	inMemoryData := NewInMemoryStorage()
	fs := FileStorage{
		inMemoryData: inMemoryData,
		filePath:     filename,
	}

	err := readFromFile(fs)
	if err != nil {
		panic(err)
	}

	return &fs
}

func readFromFile(fs FileStorage) error {
	file, err := os.OpenFile(fs.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()

		var sd ShortenedData
		if err := json.Unmarshal(line, &sd); err != nil {
			return err
		}

		err := fs.inMemoryData.Store(sd)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) Store(data ShortenedData) error {
	s.inMemoryData.data[data.ShortUrl] = data
	return nil
}

func (s *FileStorage) Get(key string) (ShortenedData, error) {
	data, exists := s.inMemoryData.data[key]
	if !exists {
		return ShortenedData{}, errors.New("key not found: " + key)
	}
	return data, nil
}

func (s *FileStorage) Save() error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, v := range s.inMemoryData.data {
		line, err := json.Marshal(v)
		if err != nil {
			return err
		}

		_, err = writer.Write(append(line, '\n'))
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
