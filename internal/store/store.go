package store

import (
	"errors"
	"fmt"
	"gg/internal/model"
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
func NewStore(data model.Data) *Store {
	s := &Store{
		Data: make(map[string]model.Address),
	}

	for k, v := range data.Storage {
		newAddress := model.Address{
			IP: v.IP,
		}

		s.Data[k] = newAddress
	}

	return s
}

func (s *Store) SearchDomain(domain string) (string, error) {
	d, ok := s.Data[domain]
	if !ok {
		return "", fmt.Errorf("domain %s not found", domain)
	}

	return d.IP, nil
}

func (s *Store) WriteToYaml(domain string, ip_string string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	if ip_string == "" {
		return nil
	}

	new_entry := model.Address{
		Domain: domain,
		IP:     ip_string,
	}

	s.Data[domain] = new_entry

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

	fmt.Printf("[Store] %s guardado en el el yaml!", domain)
	return nil
}
