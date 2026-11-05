package store

import (
	"sync"

	"github.com/Atablck/ski-rental/internal/domain"
)

type InMemoryStore struct {
	mu      sync.Mutex
	persons map[string]domain.Person
	skis    map[string]domain.Ski
	rentals map[string]domain.Rental
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		persons: make(map[string]domain.Person),
		skis:    make(map[string]domain.Ski),
		rentals: make(map[string]domain.Rental),
	}
}

func (s *InMemoryStore) AddPerson(p domain.Person) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.persons[p.ID] = p
}

func (s *InMemoryStore) GetPersons() []domain.Person {
	s.mu.Lock()
	defer s.mu.Unlock()

	var list []domain.Person
	for _, p := range s.persons {
		list = append(list, p)
	}
	return list
}

func (s *InMemoryStore) AddSki(ski domain.Ski) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.skis[ski.ID] = ski
}

func (s *InMemoryStore) GetSkis() []domain.Ski {
	s.mu.Lock()
	defer s.mu.Unlock()

	var list []domain.Ski
	for _, ski := range s.skis {
		list = append(list, ski)
	}
	return list
}

func (s *InMemoryStore) AddRental(r domain.Rental) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rentals[r.ID] = r
}

func (s *InMemoryStore) GetRentals() []domain.Rental {
	s.mu.Lock()
	defer s.mu.Unlock()

	var list []domain.Rental
	for _, r := range s.rentals {
		list = append(list, r)
	}
	return list
}