package dbmetrics

import "sync"

type Metrics struct {
	QueryCount int
	MaxTimeMs  int64
	MaxQuery   string
	mu         sync.Mutex
}

func (m *Metrics) Record(query string, durationMs int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.QueryCount++

	if durationMs > m.MaxTimeMs {
		m.MaxTimeMs = durationMs
		m.MaxQuery = query
	}
}
