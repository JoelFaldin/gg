package store

import (
	"errors"
	"fmt"
	"gg/internal/model"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(filePath string) (model.Data, error) {
	var data model.Data

	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Config file not found")
			return data, nil
		}

		return data, fmt.Errorf("failed to read file: %w", err)
	}

	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return data, fmt.Errorf("failed to parse yaml: %w", err)
	}

	return data, nil
}

// Copying the data into internal store:
func NewsStore(data model.Data) *model.Store {
	s := &model.Store{
		Data: make(map[string]model.Address),
	}

	for k, v := range data.Storage {
		newAddress := model.Address{
			Value: v.Value,
		}

		s.Data[k] = newAddress
	}

	return s
}
