package blocklist

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type BlockStruct struct {
	mu        sync.RWMutex
	domains   map[string]struct{}
	whitelist map[string]struct{}
}

func NewBlockEngine() *BlockStruct {
	return &BlockStruct{
		domains:   make(map[string]struct{}),
		whitelist: make(map[string]struct{}),
	}
}

func (b *BlockStruct) IsBlocked(domain string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	domain = strings.ToLower(strings.TrimSuffix(domain, "."))

	if _, white := b.whitelist[domain]; white {
		return false
	}

	_, black := b.domains[domain]
	return black
}

func (b *BlockStruct) LoadFromFolder(folderPath string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.domains = make(map[string]struct{})

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		f, err := os.Open(folderPath + "/" + file.Name())
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) > 1 {
				line = parts[1]
			} else {
				line = parts[0]
			}

			domain := strings.ToLower(strings.TrimSuffix(line, "."))
			b.domains[domain] = struct{}{}
		}

		f.Close()
	}

	return nil
}
