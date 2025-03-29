package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Match res la estructura que hace de base de datos
type Match struct {
	ID          int    `json:"id"`
	HomeTeam    string `json:"homeTeam"`
	AwayTeam    string `json:"awayTeam"`
	MatchDate   string `json:"matchDate"`
	GoalCount   int    `json:"goalCount,omitempty"`
	YellowCards int    `json:"yellowCards,omitempty"`
	RedCards    int    `json:"redCards,omitempty"`
	ExtraTime   bool   `json:"extraTime,omitempty"`
}

// Maneja la información de Match
type MatchRepository struct {
	mu      sync.RWMutex
	matches map[int]Match
	nextID  int
}

// Utiliza un mapa para guradar los matches
func NewMatchRepository() *MatchRepository {
	return &MatchRepository{
		matches: make(map[int]Match),
		nextID:  1,
	}
}

// GetAllMatches muestra todos los partidos
func (mr *MatchRepository) GetAllMatches() []Match {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	matches := make([]Match, 0, len(mr.matches))
	for _, match := range mr.matches {
		matches = append(matches, match)
	}
	return matches
}

// GetMatchByID muestra el partido que se busca por ID
func (mr *MatchRepository) GetMatchByID(id int) (Match, bool) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	match, exists := mr.matches[id]
	return match, exists
}

// CreateMatch crea nuevos partidos
func (mr *MatchRepository) CreateMatch(match Match) int {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	match.ID = mr.nextID
	mr.matches[match.ID] = match
	mr.nextID++
	return match.ID
}

// UpdateMatch actualiza un partido existente
func (mr *MatchRepository) UpdateMatch(id int, updatedMatch Match) bool {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, exists := mr.matches[id]; !exists {
		return false
	}

	updatedMatch.ID = id
	mr.matches[id] = updatedMatch
	return true
}

// DeleteMatch elimina un partido por su id
func (mr *MatchRepository) DeleteMatch(id int) bool {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if _, exists := mr.matches[id]; !exists {
		return false
	}

	delete(mr.matches, id)
	return true
}

// PATCH method, maneja los metodos para registrar un gol, una tarjeta amarilla, una tarjeta rojo y el timepo extra
func (mr *MatchRepository) RegisterGoal(id int) bool {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	match, exists := mr.matches[id]
	if !exists {
		return false
	}

	match.GoalCount++
	mr.matches[id] = match
	return true
}

func (mr *MatchRepository) RegisterYellowCard(id int) bool {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	match, exists := mr.matches[id]
	if !exists {
		return false
	}

	match.YellowCards++
	mr.matches[id] = match
	return true
}

func (mr *MatchRepository) RegisterRedCard(id int) bool {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	match, exists := mr.matches[id]
	if !exists {
		return false
	}

	match.RedCards++
	mr.matches[id] = match
	return true
}

func (mr *MatchRepository) SetExtraTime(id int) bool {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	match, exists := mr.matches[id]
	if !exists {
		return false
	}

	match.ExtraTime = true
	mr.matches[id] = match
	return true
}

func main() {
	router := mux.NewRouter()
	matchRepo := NewMatchRepository()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	router.HandleFunc("/api/matches", func(w http.ResponseWriter, r *http.Request) {
		matches := matchRepo.GetAllMatches()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(matches)
	}).Methods("GET")

	router.HandleFunc("/api/matches/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid match ID", http.StatusBadRequest)
			return
		}

		match, exists := matchRepo.GetMatchByID(id)
		if !exists {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
	}).Methods("GET")

	router.HandleFunc("/api/matches", func(w http.ResponseWriter, r *http.Request) {
		var match Match
		err := json.NewDecoder(r.Body).Decode(&match)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := matchRepo.CreateMatch(match)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": id})
	}).Methods("POST")

	router.HandleFunc("/api/matches/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid match ID", http.StatusBadRequest)
			return
		}

		var match Match
		err = json.NewDecoder(r.Body).Decode(&match)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !matchRepo.UpdateMatch(id, match) {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("PUT")

	router.HandleFunc("/api/matches/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid match ID", http.StatusBadRequest)
			return
		}

		if !matchRepo.DeleteMatch(id) {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	router.HandleFunc("/api/matches/{id}/goals", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid match ID", http.StatusBadRequest)
			return
		}

		if !matchRepo.RegisterGoal(id) {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("PATCH")

	router.HandleFunc("/api/matches/{id}/yellowcards", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid match ID", http.StatusBadRequest)
			return
		}

		if !matchRepo.RegisterYellowCard(id) {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("PATCH")

	router.HandleFunc("/api/matches/{id}/redcards", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid match ID", http.StatusBadRequest)
			return
		}

		if !matchRepo.RegisterRedCard(id) {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("PATCH")

	router.HandleFunc("/api/matches/{id}/extratime", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid match ID", http.StatusBadRequest)
			return
		}

		if !matchRepo.SetExtraTime(id) {
			http.Error(w, "Match not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("PATCH")

	handler := c.Handler(router)
	port := 8081
	fmt.Printf("Server is running on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
