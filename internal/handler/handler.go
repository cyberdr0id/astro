package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/cyberdr0id/astro/internal/service"
)

type messageResponse struct {
	Message string `json:"message"`
}

var (
	errInvalidDate = errors.New("invalid date")
	errEmptyAPIkey = errors.New("empty API key")
)

// GetImage handles requests for getting an images.
func (s *Server) GetImage(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("api_key")
	date := r.URL.Query().Get("date")

	err := validateRequest(date, apiKey)
	if err != nil {
		sendResponse(w, messageResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	url := "https://api.nasa.gov/planetary/apod"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		sendResponse(w, messageResponse{Message: "failed to create new request: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()

	q.Add("api_key", apiKey)
	q.Add("date", date)

	req.URL.RawQuery = q.Encode()

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		sendResponse(w, messageResponse{Message: "unable to call APOD: " + err.Error()}, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var entry service.Entry

	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		sendResponse(w, messageResponse{Message: "unable to decode response: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	if entry.MediaType != "image" {
		sendResponse(w, messageResponse{Message: "unsupported media type"}, http.StatusBadRequest)
		return
	}

	img, err := http.Get(entry.URL)
	if err != nil {
		sendResponse(w, messageResponse{Message: "unable to get image: " + err.Error()}, http.StatusInternalServerError)
		return
	}
	defer img.Body.Close()

	fileID, err := s.service.SaveImage(img.Body)
	if err != nil {
		sendResponse(w, messageResponse{Message: "unable to save image: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	id, err := s.service.SaveEntry(entry, fileID)
	if err != nil {
		sendResponse(w, messageResponse{Message: "unable to save entry: %w" + err.Error()}, http.StatusInternalServerError)
		return
	}

	sendResponse(w, messageResponse{Message: "picture has been downloaded, entry id " + id}, http.StatusOK)
}

// GetEntries handles request for getting entries.
func (s *Server) GetEntries(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	apiKey := r.URL.Query().Get("api_key")

	err := validateRequest(date, apiKey)
	if err != nil {
		sendResponse(w, messageResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	entries, err := s.service.GetEntries(date)
	if err != nil {
		sendResponse(w, messageResponse{Message: "unable to retrieve entries from the database: " + err.Error()}, http.StatusBadRequest)
		return
	}

	sendResponse(w, entries, http.StatusOK)
}

func validateRequest(date, apiKey string) error {
	dateExp := "^\\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$"

	if apiKey == "" {
		return errEmptyAPIkey
	}

	ok, err := regexp.MatchString(dateExp, date)
	if !ok && date != "" {
		return errInvalidDate
	}
	if err != nil {
		return fmt.Errorf("cannot validate date: %w", err)
	}

	return nil
}

func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
