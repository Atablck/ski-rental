package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Atablck/ski-rental/internal/domain"
	"github.com/Atablck/ski-rental/internal/handler"
	"github.com/Atablck/ski-rental/internal/store"
)

func main() {
	err := os.MkdirAll("./static/uploads", 0755)
	if err != nil {
		fmt.Printf("Fehler beim Erstellen des Upload-Ordners: %v\n", err)
	}

	db := store.NewStore()

	db.AddPerson(domain.Person{ID: "1", FirstName: "Max", LastName: "Mustermann", AusweisNr: "AT12345", AusweisTyp: "Reisepass"})

	db.AddItem(domain.Ski{ID: "s1", Name: "Atomic Redster", Length: 170, PricePerHour: 12.5})
	db.AddItem(domain.Ski{ID: "s2", Name: "Atomic Redster S9", Length: 165, PricePerHour: 15.0})
	db.AddItem(domain.Ski{ID: "s3", Name: "Fischer RC4", Length: 170, PricePerHour: 12.0})
	db.AddItem(domain.Ski{ID: "s4", Name: "Head Supershape", Length: 163, PricePerHour: 14.5})
	db.AddItem(domain.Ski{ID: "s5", Name: "Rossignol Hero", Length: 180, PricePerHour: 13.0})

	db.AddItem(domain.Bike{ID: "b1", Name: "Mountainbike Trek", FrameSize: "L", PricePerHour: 8.0})
	db.AddItem(domain.Bike{ID: "b2", Name: "Cube Attention", FrameSize: "M", PricePerHour: 9.5})
	db.AddItem(domain.Bike{ID: "b3", Name: "Specialized Rockhopper", FrameSize: "S", PricePerHour: 11.0})
	db.AddItem(domain.Bike{ID: "b4", Name: "Cannondale Trail", FrameSize: "XL", PricePerHour: 10.0})
	db.AddItem(domain.Bike{ID: "b5", Name: "Giant Talon", FrameSize: "M", PricePerHour: 8.5})

	h := &handler.Handler{Store: db}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/stats", h.GetStats)
	mux.HandleFunc("/api/persons", h.CreatePerson)
	mux.HandleFunc("/api/persons/list", h.ListPersons)
	mux.HandleFunc("/api/items/available", h.GetAvailableItems)
	mux.HandleFunc("/api/rentals", h.CreateRental)
	mux.HandleFunc("/api/items/all", h.GetAllItems)
	mux.HandleFunc("/api/rentals/active", h.GetActiveRentals)
	mux.HandleFunc("/api/rentals/return", h.ReturnItem)

	mux.Handle("/", http.FileServer(http.Dir("./static")))

	port := ":8080"
	fmt.Printf("Server startet auf http://localhost%s...\n", port)

	serveErr := http.ListenAndServe(port, mux)
	if serveErr != nil {
		fmt.Printf("Fehler beim Starten des Servers: %v\n", serveErr)
	}
}