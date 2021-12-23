package store

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/vineet-garg/assignment-1/config"
	"sync"
	"sync/atomic"
	"time"
)

// Store supports aadding ang getting of paassword Hash
type Store interface {
	// AddPwd returns a new id and adds password hash after a defay of preconfigured duration
	AddPwd(pwd []byte, wg *sync.WaitGroup) int64
	GetHash(id int64) (string, bool)
}

// GetStore returns the interface to the internal store
func GetStore() Store {
	return &internal
}



// Singleton vale of internal store
var internal = pwdStore{}

// Internal Types
type pwdHash struct {
	salt []byte
	algo string
	hash string
}

type pwdStore struct {
	id      int64
	safeMap sync.Map
	Store
}

func (p *pwdStore) AddPwd(pwd []byte, wg *sync.WaitGroup) int64 {
	newID := atomic.AddInt64(&p.id, 1)
	h := sha512.Sum512(pwd)
	fmt.Println()
	item := pwdHash{salt: []byte{},
		algo: "SHA-512",
		hash: base64.StdEncoding.EncodeToString(h[:]),
	}
	go func() {
		timer1 := time.NewTimer(config.Delay)
		<-timer1.C
		p.safeMap.Store(newID, item)
		wg.Done()
	}()
	return newID
}

func (p *pwdStore) GetHash(id int64) (string, bool) {
	if result, ok := p.safeMap.Load(id); ok {
		return result.(pwdHash).hash, ok
	} else {
		return "", false
	}
}
