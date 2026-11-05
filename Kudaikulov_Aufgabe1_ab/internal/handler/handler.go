package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Atablck/ski-rental/internal/domain"
	"github.com/Atablck/ski-rental/internal/store"
)

type Handler struct {
	Store *store.InMemoryStore
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	persons := h.Store.GetPersons()
	skis := h.Store.GetSkis()
	rentals := h.Store.GetRentals()

	activeRentals := 0
	for _, rent := range rentals {
		if rent.IsActive {
			activeRentals++
		}
	}

	stats := map[string]int{
		"total_items":     len(skis),
		"rented_items":    activeRentals,
		"available_items": len(skis) - activeRentals,
		"total_persons":   len(persons),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST erlaubt", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Fehler beim Parsen: %v", err), http.StatusBadRequest)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	ausweisNr := r.FormValue("ausweis_nr")
	ausweisTyp := r.FormValue("ausweis_typ")

	var ausweisPfad string
	file, fileHeader, err := r.FormFile("ausweis_file")
	if err == nil {
		defer file.Close()
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(fileHeader.Filename))
		filePath := filepath.Join(".", "static", "uploads", fileName)
		
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Fehler beim Speichern der Datei: %v", err), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		io.Copy(dst, file)
		
		ausweisPfad = "/uploads/" + fileName 
	} else if err != http.ErrMissingFile {
		http.Error(w, fmt.Sprintf("Fehler beim Datei-Upload: %v", err), http.StatusBadRequest)
		return
	}

	persons := h.Store.GetPersons()
	newID := fmt.Sprintf("p%d", len(persons)+1)

	p := domain.Person{
		ID:          newID,
		FirstName:   firstName,
		LastName:    lastName,
		AusweisNr:   ausweisNr,
		AusweisTyp:  ausweisTyp,
		AusweisPfad: ausweisPfad,
	}
	h.Store.AddPerson(p)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "erfolgreich"}`))
}

func (h *Handler) ListPersons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Nur GET erlaubt", http.StatusMethodNotAllowed)
		return
	}

	persons := h.Store.GetPersons()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

func (h *Handler) GetAvailableItems(w http.ResponseWriter, r *http.Request) {
	skis := h.Store.GetSkis()
	rentals := h.Store.GetRentals()

	activeItemIDs := make(map[string]bool)
	for _, rent := range rentals {
		if rent.IsActive {
			activeItemIDs[rent.ItemID] = true
		}
	}

	var available []domain.Ski
	for _, ski := range skis {
		if !activeItemIDs[ski.ID] {
			available = append(available, ski)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(available)
}

func (h *Handler) CreateRental(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST erlaubt", http.StatusMethodNotAllowed)
		return
	}

	personID := r.FormValue("person_id")
	itemID := r.FormValue("item_id")
	hoursStr := r.FormValue("hours")

	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		hours = 1 
	}

	rentals := h.Store.GetRentals()
	newID := fmt.Sprintf("r%d", len(rentals)+1)

	rent := domain.Rental{
		ID:       newID,
		PersonID: personID,
		ItemID:   itemID,
		Start:    time.Now(),
		Hours:    hours,
		IsActive: true,
	}

	h.Store.AddRental(rent)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "erfolgreich"}`))
}

func (h *Handler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	skis := h.Store.GetSkis()

	type ItemResponse struct {
		ID    string  `json:"id"`
		Name  string  `json:"name"`
		Type  string  `json:"type"`
		Price float64 `json:"price"`
	}

	var response []ItemResponse
	for _, ski := range skis {
		response = append(response, ItemResponse{
			ID:    ski.ID,
			Name:  ski.Name,
			Type:  "Ski",
			Price: ski.PricePerHour,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetActiveRentals(w http.ResponseWriter, r *http.Request) {
	rentals := h.Store.GetRentals()
	persons := h.Store.GetPersons()
	skis := h.Store.GetSkis()

	personNames := make(map[string]string)
	for _, p := range persons {
		personNames[p.ID] = p.FirstName + " " + p.LastName
	}

	itemNames := make(map[string]string)
	for _, s := range skis {
		itemNames[s.ID] = s.Name
	}

	type ActiveResponse struct {
		RentalID   string `json:"rental_id"`
		PersonName string `json:"person_name"`
		ItemName   string `json:"item_name"`
		Hours      int    `json:"hours"`
	}

	var response []ActiveResponse
	for _, rent := range rentals {
		if rent.IsActive {
			response = append(response, ActiveResponse{
				RentalID:   rent.ID,
				PersonName: personNames[rent.PersonID],
				ItemName:   itemNames[rent.ItemID],
				Hours:      rent.Hours,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ReturnItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST erlaubt", http.StatusMethodNotAllowed)
		return
	}

	rentalID := r.FormValue("rental_id")
	rentals := h.Store.GetRentals()
	skis := h.Store.GetSkis()

	var rent domain.Rental
	var found bool
	for _, req := range rentals {
		if req.ID == rentalID {
			rent = req
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Ausleihe nicht gefunden", http.StatusNotFound)
		return
	}

	var rate float64
	for _, s := range skis {
		if s.ID == rent.ItemID {
			rate = s.PricePerHour
			break
		}
	}

	finalPrice := rent.CalcPrice(rate, time.Now())

	rent.IsActive = false
	h.Store.AddRental(rent)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"price": finalPrice})
}