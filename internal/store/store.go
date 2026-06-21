package store

import (
	"errors"
	"fmt"
	"gg/internal/model"
	"log/slog"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Store struct {
	Data map[string]model.MemoryAddress
}

var fileMutex sync.Mutex

func LoadConfig(filePath string) (model.MemoryData, error) {
	data := model.MemoryData{
		Storage: make(map[string]model.MemoryAddress),
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Warn("Data file not found, make sure .env has a valid path to yaml")
			return data, nil
		}

		return data, fmt.Errorf("failed to read file: %w", err)
	}

	var yamlData struct {
		Storage map[string]model.Address `yaml:"storage"`
	}

	err = yaml.Unmarshal(file, &yamlData)
	if err != nil {
		return data, fmt.Errorf("failed to parse yaml: %w", err)
	}

	for domain, addr := range yamlData.Storage {
		exp := time.Now().Add(time.Duration(addr.TTL) * time.Second)

		data.Storage[domain] = model.MemoryAddress{
			TTL:  exp,
			IPv4: addr.IPv4,
			IPv6: addr.IPv6,
		}
	}

	return data, nil
}

// Copying the data into internal store:
func NewStore(data model.MemoryData) *Store {
	s := &Store{
		Data: make(map[string]model.MemoryAddress),
	}

	for k, v := range data.Storage {
		newAddress := model.MemoryAddress{
			IPv4: v.IPv4,
			IPv6: v.IPv6,
			TTL:  v.TTL,
		}

		s.Data[k] = newAddress
	}

	return s
}

func (s *Store) SearchDomain(domain string) (model.MemoryAddress, bool) {
	d, ok := s.Data[domain]
	if !ok {
		return model.MemoryAddress{}, false
	}

	return d, true
}

func (s *Store) WriteToYaml(domain string, ips []string, questionType uint16, ttl int) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	if len(ips) == 0 {
		return nil
	}

	// Check if domain already exists:
	record, exists := s.Data[domain]
	if !exists {
		record = model.MemoryAddress{}
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

	record.TTL = time.Now().Add(time.Duration(ttl) * time.Second)
	s.Data[domain] = record

	// Convert all s.Data (MemoryAddress) to model.Address
	yamlStorage := make(map[string]model.Address)

	for dom, memAddr := range s.Data {
		ttlRemainder := int(time.Until(memAddr.TTL).Seconds())

		if ttlRemainder < 0 {
			ttlRemainder = 0
		}

		yamlStorage[dom] = model.Address{
			TTL:  ttlRemainder,
			IPv4: memAddr.IPv4,
			IPv6: memAddr.IPv6,
		}
	}

	data_to_save := model.Data{
		Storage: yamlStorage,
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
