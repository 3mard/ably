package server

import (
	"fmt"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var (
	ErrorNotFound = fmt.Errorf("not found")
)

type ClientRepository interface {
	GetClientMessages(client_id string) ([]int32, error)
	SetClientMessages(client_id string, message_id []int32) error
}

type InMemoryClientRepository struct {
	ttl                time.Duration
	clientMessageStore *ttlcache.Cache[string, []int32]
}

func NewInMemoryClientRepository(ttl time.Duration) *InMemoryClientRepository {
	return &InMemoryClientRepository{
		ttl:                ttl,
		clientMessageStore: ttlcache.New(ttlcache.WithTTL[string, []int32](ttl)),
	}
}

func (r *InMemoryClientRepository) GetClientMessages(client_id string) ([]int32, error) {
	item := r.clientMessageStore.Get(client_id)
	if item == nil {
		return nil, ErrorNotFound
	}
	if item.IsExpired() {
		return nil, ErrorNotFound
	}

	return item.Value(), nil
}

func (r *InMemoryClientRepository) SetClientMessages(client_id string, message_id []int32) error {

	r.clientMessageStore.Set(client_id, message_id, r.ttl)
	return nil
}
