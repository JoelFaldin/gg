package store

import (
	"errors"
	"fmt"
	"gg/internal/model"
	"log/slog"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Store struct {
	Data map[string]model.Address
}

var fileMutex sync.Mutex

func LoadConfig(filePath string) (model.Data, error) {
	var data model.Data

	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Warn("Data file not found, make sure .env has a valid path to yaml")
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
func NewStore(data model.Data) *Store {
	s := &Store{
		Data: make(map[string]model.Address),
	}

	for k, v := range data.Storage {
		newAddress := model.Address{
			IPv4: v.IPv4,
			IPv6: v.IPv6,
		}

		s.Data[k] = newAddress
	}

	return s
}

func (s *Store) SearchDomain(domain string) (model.Address, bool) {
	d, ok := s.Data[domain]
	if !ok {
		return model.Address{}, false
	}

	return d, true
}

func (s *Store) WriteToYaml(domain string, ips []string, questionType uint16) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	if len(ips) == 0 {
		return nil
	}

	// Check if domain already exists:
	record, exists := s.Data[domain]
	if !exists {
		record = model.Address{}
	}

	switch questionType {
	case 1:
		record.IPv4 = ips
	case 28:
		record.IPv6 = ips
	default:
		slog.Warn("Resource not supported (yet!)")
		return nil
	}

	s.Data[domain] = record

	data_to_save := model.Data{
		Storage: s.Data,
	}

	res, err := yaml.Marshal(data_to_save)
	if err != nil {
		return err
	}

	err = os.WriteFile("data.yaml", res, 0644)
	if err != nil {
		return err
	}

	slog.Debug(fmt.Sprintf("[Store] %s saved on YAML!", domain))

	return nil
}
