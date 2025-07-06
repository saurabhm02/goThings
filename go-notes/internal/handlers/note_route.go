package handlers

import (
	"encoding/json"
	"go-notes/internal/models"
	"go-notes/internal/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type NoteHandler struct {
	Service *services.NoteService
}

func NewNoteHandler(service *services.NoteService) *NoteHandler {
	return &NoteHandler{
		Service: service,
	}
}

func (s *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		resp := models.Response{
			Status:  "error",
			Message: "Invalid request body for creating note",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp, err := s.Service.CreateNote(note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	resp, err := s.Service.GetAllNotes()
	if err != nil {
		res := models.Response{
			Status:  "error",
			Message: "error while creating note",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *NoteHandler) GetNoteById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteIdStr := vars["id"]
	noteID, err := strconv.Atoi(noteIdStr)
	if err != nil {
		resp := models.Response{
			Status:  "error",
			Message: "Invalid note ID in URL",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp, err := s.Service.GetNoteById(noteID)
	if err != nil {
		res := models.Response{
			Status:  "error",
			Message: "Error while getting note by ID",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	noteIdStr := vars["id"]
	noteID, err := strconv.Atoi(noteIdStr)
	if err != nil {
		res := models.Response{
			Status:  "error",
			Message: "Invalid note id",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		res := models.Response{
			Status:  "error",
			Message: "Invalid request body for updating note",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	resp, err := s.Service.UpdateNote(noteID, note)
	if err != nil {
		res := models.Response{
			Status:  "error",
			Message: "error while updaing note by id",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteIdStr := vars["id"]

	noteId, err := strconv.Atoi(noteIdStr)
	if err != nil {
		resp := models.Response{
			Status:  "error",
			Message: "Invalid note ID in URL",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp, err := s.Service.DeleteNote(noteId)
	if err != nil {
		res := models.Response{
			Status:  "error",
			Message: "error while deleting note",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
