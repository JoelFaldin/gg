package blocklist

import (
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
