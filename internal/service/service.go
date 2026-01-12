package service

import (
	"assignment2/internal/storage"
	"sync/atomic"
)

type DataService struct {
	storage      *storage.DataStorage
	totalRequest int64
}

func NewDataService(storage *storage.DataStorage) *DataService {
	return &DataService{
		storage: storage,
	}
}

func (s *DataService) SaveKeyValue(key, value string) {
	atomic.AddInt64(&s.totalRequest, 1)
	s.storage.Set(key, value)
}

func (s *DataService) GetAll() map[string]string {
	atomic.AddInt64(&s.totalRequest, 1)
	return s.storage.GetAll()
}

func (s *DataService) DeleteKey(key string) bool {
	atomic.AddInt64(&s.totalRequest, 1)
	return s.storage.Delete(key)
}

func (s *DataService) GetStats() (int64, int) {
	atomic.AddInt64(&s.totalRequest, 1)
	return atomic.LoadInt64(&s.totalRequest), s.storage.Size()
}

func (s *DataService) GetCurrentStats() (int64, int) {
	return atomic.LoadInt64(&s.totalRequest), s.storage.Size()
}
