package helper

import (
	"time"
)

type metric struct {
	startTime time.Time
}

func NewMetric() *metric {
	return &metric{
		startTime: time.Now(),
	}
}

func (m *metric) Reset() {
	m.startTime = time.Now()
}

func (m *metric) Elapsed() time.Duration {
	return time.Since(m.startTime)
}

func (m *metric) ElapsedAndReset() time.Duration {
	elapsed := m.Elapsed()
	m.Reset()
	return elapsed
}
