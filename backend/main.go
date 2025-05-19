package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Item struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

var dataFile = "data.json"

// Load items from the file
func loadItems() ([]Item, error) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	var items []Item
	if len(data) > 0 {
		if err := json.Unmarshal(data, &items); err != nil {
			return nil, err
		}
	}
	return items, nil
}

// Save items to the file
func saveItems(items []Item) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

// Handlers
func getAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := loadItems()
	if err != nil {
		http.Error(w, "Failed to read items", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	items, err := loadItems()
	if err != nil {
		http.Error(w, "Failed to load items", http.StatusInternalServerError)
		return
	}

	// Set ID
	newItem.ID = 1
	if len(items) > 0 {
		newItem.ID = items[len(items)-1].ID + 1
	}
	items = append(items, newItem)

	if err := saveItems(items); err != nil {
		http.Error(w, "Failed to save item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	updatedItem.ID = id

	items, err := loadItems()
	if err != nil {
		http.Error(w, "Failed to load items", http.StatusInternalServerError)
		return
	}

	updated := false
	for i, item := range items {
		if item.ID == id {
			items[i] = updatedItem
			updated = true
			break
		}
	}
	if !updated {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if err := saveItems(items); err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedItem)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	items, err := loadItems()
	if err != nil {
		http.Error(w, "Failed to load items", http.StatusInternalServerError)
		return
	}

	newItems := []Item{}
	found := false
	for _, item := range items {
		if item.ID != id {
			newItems = append(newItems, item)
		} else {
			found = true
		}
	}

	if !found {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if err := saveItems(newItems); err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, //frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Routes
	r.Route("/api/events", func(r chi.Router) {
		r.Get("/", getAllItems)       // GET /api/events
		r.Post("/", createItem)       // POST /api/events
		r.Put("/{id}", updateItem)    // PUT /api/events/123
		r.Delete("/{id}", deleteItem) // DELETE /api/events/123
	})

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
