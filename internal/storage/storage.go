package storage

import "sync"

type DataStorage struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewDataStorage() *DataStorage {
	return &DataStorage{
		data: make(map[string]string),
	}
}

func (s *DataStorage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *DataStorage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *DataStorage) GetAll() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]string, len(s.data))
	for k, v := range s.data {
		result[k] = v
	}
	return result
}

func (s *DataStorage) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return true
	}
	return false
}

func (s *DataStorage) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
