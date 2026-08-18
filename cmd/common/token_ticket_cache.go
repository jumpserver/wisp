package common

import (
	"sync"
	"time"
)

const (
	tokenTicketCacheMaxEntries = 10000
	tokenTicketCacheMaxTTL     = 24 * time.Hour
)

type tokenTicketEntry struct {
	ticketID string
	expires  int64
}

type TokenTicketCache struct {
	mu         sync.Mutex
	items      map[string]tokenTicketEntry
	maxEntries int
}

func NewTokenTicketCache() *TokenTicketCache {
	return newTokenTicketCache(tokenTicketCacheMaxEntries)
}

func newTokenTicketCache(maxEntries int) *TokenTicketCache {
	if maxEntries < 1 {
		maxEntries = 1
	}
	return &TokenTicketCache{
		items:      make(map[string]tokenTicketEntry),
		maxEntries: maxEntries,
	}
}

func (c *TokenTicketCache) Store(tokenID, ticketID string, expires int64) {
	now := time.Now().Unix()
	if tokenID == "" {
		return
	}
	if ticketID == "" || expires <= now {
		c.Delete(tokenID)
		return
	}
	maxExpires := now + int64(tokenTicketCacheMaxTTL/time.Second)
	if expires > maxExpires {
		expires = maxExpires
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.items[tokenID]; !ok && len(c.items) >= c.maxEntries {
		c.evict(now)
	}
	c.items[tokenID] = tokenTicketEntry{
		ticketID: ticketID,
		expires:  expires,
	}
}

func (c *TokenTicketCache) Take(tokenID string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[tokenID]
	if !ok {
		return "", false
	}
	delete(c.items, tokenID)
	if item.expires <= time.Now().Unix() {
		return "", false
	}
	return item.ticketID, true
}

func (c *TokenTicketCache) Delete(tokenID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, tokenID)
}

func (c *TokenTicketCache) evict(now int64) {
	var earliestKey string
	var earliestExpires int64
	for tokenID, item := range c.items {
		if item.expires <= now {
			delete(c.items, tokenID)
			continue
		}
		if earliestKey == "" || item.expires < earliestExpires {
			earliestKey = tokenID
			earliestExpires = item.expires
		}
	}
	if len(c.items) >= c.maxEntries && earliestKey != "" {
		delete(c.items, earliestKey)
	}
}
