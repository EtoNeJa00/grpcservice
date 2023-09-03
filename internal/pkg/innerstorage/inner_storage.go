package innerstorage

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("not found")
)

type InnerStorage interface {
	Get(id uuid.UUID) (string, error)
	Create(record string) (uuid.UUID, string)
	Update(id uuid.UUID, record string) (string, error)
	Delete(id uuid.UUID) (string, error)

	CleanUp()
}

type innerStorage struct {
	storage map[uuid.UUID]string
	mu      sync.RWMutex
}

func NewInnerStorage() InnerStorage {
	return &innerStorage{
		storage: map[uuid.UUID]string{},
	}
}

func (s *innerStorage) Get(id uuid.UUID) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.storage[id]
	if !ok {
		return "", ErrNotFound
	}

	return v, nil
}

func (s *innerStorage) Create(record string) (id uuid.UUID, recordResult string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id = uuid.New()
	s.storage[id] = record

	return id, record
}

func (s *innerStorage) Update(id uuid.UUID, record string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.storage[id]
	if !ok {
		return "", ErrNotFound
	}

	s.storage[id] = record

	return record, nil
}

func (s *innerStorage) Delete(id uuid.UUID) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.storage[id]
	if !ok {
		return "", ErrNotFound
	}

	delete(s.storage, id)

	return v, nil
}

func (s *innerStorage) CleanUp() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage = map[uuid.UUID]string{}
}
