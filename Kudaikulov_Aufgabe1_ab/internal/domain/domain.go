package domain

import (
	"math"
	"time"
)

type Person struct {
	ID          string
	FirstName   string
	LastName    string
	AusweisNr   string
	AusweisTyp  string
	AusweisPfad string 
}

type Ski struct {
	ID           string
	Name         string
	Length       int
	PricePerHour float64
}

type Rental struct {
	ID       string
	PersonID string
	ItemID   string
	Start    time.Time
	Hours    int 
	IsActive bool
}

func (r *Rental) CalcPrice(rate float64, returnTime time.Time) float64 {
	duration := returnTime.Sub(r.Start).Hours()
	hours := int(math.Ceil(duration))
	
	if hours < 1 {
		hours = 1
	}
	return float64(hours) * rate
}