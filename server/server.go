package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cfs/db"
	"cfs/proc"
)

type Server struct {
	db         db.Database
	classifier proc.Classifier
}

func (s *Server) Init() error {
	s.db = db.Database{}
	if err := s.db.Init(); err != nil {
		return err
	}
	s.db.Seed()
	s.classifier = proc.Classifier{}
	categories, err := s.db.GetCategories()
	if err != nil {
		return err
	}
	s.classifier.Init(categories)
	return nil
}

func (s *Server) Close() {
	s.db.Close()
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /cfs/i", s.handleGetClassifications)
	mux.HandleFunc("POST /cfs/i", s.handleCreateClassifications)
	mux.HandleFunc("GET /cfs/i/{item}", s.handleGetClassification)
	mux.HandleFunc("GET /cfs/c", s.handleGetCategories)
	mux.HandleFunc("POST /cfs/c", s.handleCreateCategories)
	mux.HandleFunc("GET /cfs/c/{category}", s.handleGetCategory)

	fmt.Println("Server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (s *Server) handleGetClassifications(w http.ResponseWriter, r *http.Request) {
	classifications, err := s.db.GetClassifications()
	if err != nil {
		http.Error(w, "Failed to get classifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proc.ClassificationOutputData{Results: classifications})
}

func (s *Server) handleGetClassification(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	if item == "" {
		http.Error(w, "Missing item", http.StatusBadRequest)
		return
	}

	result, err := s.db.GetClassification(item)
	if err != nil {
		http.Error(w, "Failed to get classification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleCreateClassifications(w http.ResponseWriter, r *http.Request) {
	var inputData proc.InputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	results := make([]proc.ClassificationResult, 0)
	for _, item := range inputData.Items {
		classification := s.classifier.Classify(item)
		result := proc.ClassificationResult{
			Item:       item,
			Category:   classification.Category,
			Confidence: classification.Confidence,
			Matches:    classification.Matches,
		}
		if err := s.db.AddClassification(item, result); err != nil {
			http.Error(w, "Failed to create classification", http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}

	response := proc.ClassificationOutputData{Results: results}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := s.db.GetCategories()
	if err != nil {
		http.Error(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proc.CategoryOutputData{Categories: categories})
}

func (s *Server) handleCreateCategories(w http.ResponseWriter, r *http.Request) {
	var categories []proc.Category
	if err := json.NewDecoder(r.Body).Decode(&categories); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	for _, category := range categories {
		if err := s.db.AddCategory(category); err != nil {
			http.Error(w, "Failed to create category", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proc.CategoryOutputData{Categories: categories})
}

func (s *Server) handleGetCategory(w http.ResponseWriter, r *http.Request) {
	categoryName := r.URL.Query().Get("category")
	if categoryName == "" {
		http.Error(w, "Missing category", http.StatusBadRequest)
		return
	}

	category, err := s.db.GetCategory(categoryName)
	if err != nil {
		http.Error(w, "Failed to get category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}
