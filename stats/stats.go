// stats exposes interface and singelton store where stats can be stored and read from.
package stats

import (
	"sync"
	"time"
)
// Stats supports updating and getting of stats
type Stats interface {
	Update(d time.Duration)
	Get() (int64, int64)
}

// GetStatsStore returns the interface to the internal stats store
func GetStatsStore() Stats {
	return &internal
}



// Singleton vale of internal store
var internal = statsStore{}

// Internal Types and implementations
type statsStore struct {
	sync.RWMutex
	count   int64
	latency time.Duration
	Stats
}

func (s *statsStore) Update(d time.Duration) {
	s.Lock()
	defer s.Unlock()
	s.count += +1
	s.latency += d
}

func (s *statsStore) Get() (int64, int64) {
	s.RLock()
	defer s.RUnlock()
	micro := s.latency.Microseconds()
	if s.count == 0 {
		return 0, 0
	}
	return s.count, micro / int64(s.count)
}