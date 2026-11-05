package store

import (
	"sync"
	"github.com/Atablck/ski-rental/internal/domain"
)

type InMemoryStore struct {
	mu      sync.Mutex 
	persons map[string]domain.Person
	items   map[string]domain.RentalItem
	rentals map[string]domain.Rental
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		persons: make(map[string]domain.Person),
		items:   make(map[string]domain.RentalItem),
		rentals: make(map[string]domain.Rental),
	}
}


// Person Methods 

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


// Item Methods 

func (s *InMemoryStore) AddItem(item domain.RentalItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item.GetID()] = item
}

func (s *InMemoryStore) GetItems() []domain.RentalItem {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	var list []domain.RentalItem
	for _, item := range s.items {
		list = append(list, item)
	}
	return list
}


// Rental Methods

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