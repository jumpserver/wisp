package common

import (
	"testing"
	"time"
)

func TestTokenTicketCacheTake(t *testing.T) {
	cache := newTokenTicketCache(2)
	cache.Store("token-1", "ticket-1", time.Now().Add(time.Minute).Unix())

	ticketID, ok := cache.Take("token-1")
	if !ok || ticketID != "ticket-1" {
		t.Fatalf("unexpected ticket: %q, %t", ticketID, ok)
	}
	if _, ok = cache.Take("token-1"); ok {
		t.Fatal("ticket was not consumed")
	}
}

func TestTokenTicketCacheBounded(t *testing.T) {
	cache := newTokenTicketCache(2)
	now := time.Now()
	cache.Store("expired", "ticket-expired", now.Add(-time.Minute).Unix())
	if _, ok := cache.Take("expired"); ok {
		t.Fatal("expired entry was stored")
	}
	cache.Store("token-1", "ticket-1", now.Add(time.Minute).Unix())
	cache.Store("token-2", "ticket-2", now.Add(2*time.Minute).Unix())
	cache.Store("token-3", "ticket-3", now.Add(3*time.Minute).Unix())

	if len(cache.items) != 2 {
		t.Fatalf("cache size %d exceeds limit", len(cache.items))
	}
	if _, ok := cache.Take("token-1"); ok {
		t.Fatal("earliest entry was not evicted")
	}
}

func TestTokenTicketCacheMaxTTL(t *testing.T) {
	cache := newTokenTicketCache(1)
	maxExpiresBeforeStore := time.Now().Add(tokenTicketCacheMaxTTL).Unix()
	cache.Store("token-1", "ticket-1", time.Now().Add(24*time.Hour).Unix())
	maxExpiresAfterStore := time.Now().Add(tokenTicketCacheMaxTTL).Unix()

	item := cache.items["token-1"]
	if item.expires < maxExpiresBeforeStore || item.expires > maxExpiresAfterStore {
		t.Fatalf("cache expiration %d exceeds max TTL", item.expires)
	}
}
