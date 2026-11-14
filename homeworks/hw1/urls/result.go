package urls

import (
	"maps"
	"sync"
)

type SimpleAggregatorResult struct {
	mu      sync.RWMutex
	counter map[int]int
	done    bool
}

func NewSimpleAggregatorResult() *SimpleAggregatorResult {
	return &SimpleAggregatorResult{
		counter: make(map[int]int),
		done:    false,
	}
}

func (s *SimpleAggregatorResult) GetResponsesCount(code int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.counter[code]
}

func (s *SimpleAggregatorResult) count(code int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter[code]++
}

func (s *SimpleAggregatorResult) GetResult() map[int]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var curStateCounter map[int]int
	curStateCounter = maps.Clone(s.counter)

	return curStateCounter
}

func (s *SimpleAggregatorResult) finish() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.done = true
}

func (s *SimpleAggregatorResult) Done() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.done
}
