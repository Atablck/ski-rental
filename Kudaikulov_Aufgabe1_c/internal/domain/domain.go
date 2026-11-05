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

type RentalItem interface {
	GetID() string
	GetName() string
	GetType() string
	GetPricePerHour() float64
}

type Ski struct {
	ID           string
	Name         string
	Length       int
	PricePerHour float64
}

func (s Ski) GetID() string            { return s.ID }
func (s Ski) GetName() string          { return s.Name }
func (s Ski) GetType() string          { return "Ski" }
func (s Ski) GetPricePerHour() float64 { return s.PricePerHour }

type Bike struct {
	ID           string
	Name         string
	FrameSize    string
	PricePerHour float64
}

func (b Bike) GetID() string            { return b.ID }
func (b Bike) GetName() string          { return b.Name }
func (b Bike) GetType() string          { return "Bike" }
func (b Bike) GetPricePerHour() float64 { return b.PricePerHour }

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