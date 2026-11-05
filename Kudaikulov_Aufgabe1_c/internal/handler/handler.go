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
	items := h.Store.GetItems()
	rentals := h.Store.GetRentals()

	activeRentals := 0
	for _, rent := range rentals {
		if rent.IsActive {
			activeRentals++
		}
	}

	stats := map[string]int{
		"total_items":     len(items),
		"rented_items":    activeRentals,
		"available_items": len(items) - activeRentals,
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
	items := h.Store.GetItems()
	rentals := h.Store.GetRentals()

	activeItemIDs := make(map[string]bool)
	for _, rent := range rentals {
		if rent.IsActive {
			activeItemIDs[rent.ItemID] = true
		}
	}

	var available []domain.RentalItem
	for _, item := range items {
		if !activeItemIDs[item.GetID()] {
			available = append(available, item)
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
	items := h.Store.GetItems()

	type ItemResponse struct {
		ID    string  `json:"id"`
		Name  string  `json:"name"`
		Type  string  `json:"type"`
		Price float64 `json:"price"`
	}

	var response []ItemResponse
	for _, item := range items {
		response = append(response, ItemResponse{
			ID:    item.GetID(),
			Name:  item.GetName(),
			Type:  item.GetType(),
			Price: item.GetPricePerHour(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetActiveRentals(w http.ResponseWriter, r *http.Request) {
	rentals := h.Store.GetRentals()
	persons := h.Store.GetPersons()
	items := h.Store.GetItems()

	personNames := make(map[string]string)
	for _, p := range persons {
		personNames[p.ID] = p.FirstName + " " + p.LastName
	}

	itemNames := make(map[string]string)
	for _, i := range items {
		itemNames[i.GetID()] = i.GetName()
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
	items := h.Store.GetItems()

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
	for _, item := range items {
		if item.GetID() == rent.ItemID {
			rate = item.GetPricePerHour()
			break
		}
	}

	
	finalPrice := rent.CalcPrice(rate, time.Now())

	rent.IsActive = false
	h.Store.AddRental(rent)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"price": finalPrice})
}