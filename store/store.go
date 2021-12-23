// store exposes interfaces and singleton store where hashes are stored and read.
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

// Store supports adding and getting of password Hash
type Store interface {
	// AddPwd returns a new id and adds password hash after a defay of preconfigured duration
	AddPwd(pwd []byte, wg *sync.WaitGroup) int64
	GetHash(id int64) (string, bool)
}

// GetStore returns the interface to the internal store
func GetStore() Store {
	return &internal
}

// Internal Types and implementations and Values


// Singleton value of internal store
var internal = pwdStore{}

// Internal type (subject to change)
type pwdHash struct {
	// TODO salt is not the current requirement, but a cheap way to mitigate Rainbow table attacks
	salt []byte
	// storing algo along with the data helps in move to a new algo in phases.
	algo string
	hash string
}

type pwdStore struct {
	// Assumtion count will not go beyond 2^63 -1
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
