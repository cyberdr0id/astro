package service

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/cyberdr0id/astro/internal/repository"
	"github.com/cyberdr0id/astro/internal/storage"
	"github.com/pborman/uuid"
)

// Entry presents a type for receiving data from NASA API.
type Entry struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	HDURL          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

// APOD an interface for manipulating with images and entries.
type APOD interface {
	SaveImage(file io.ReadCloser, entryID string) error
	GetEntries(date string) ([]repository.Entry, error)
	SaveEntry(entry Entry) (string, error)
}

// PhotoService a type that implements APOD interface.
type PhotoService struct {
	repo *repository.Repository
	s3   *storage.S3
}

// New creates a new instance of PhotoService.
func New(repo *repository.Repository, storage *storage.S3) *PhotoService {
	return &PhotoService{
		repo: repo,
		s3:   storage,
	}
}

// SaveEntry does all actions for adding entry to the database.
func (s *PhotoService) SaveEntry(e Entry) (string, error) {
	id, err := s.repo.SaveEntry(
		e.Copyright,
		e.Date,
		e.Explanation,
		e.HDURL,
		e.MediaType,
		e.ServiceVersion,
		e.Title,
		e.URL,
	)
	if err != nil {
		return "", fmt.Errorf("unable to save entry: %w", err)
	}

	return id, nil
}

// SaveImage handle input file and add it to the database.
func (s *PhotoService) SaveImage(file io.ReadCloser, entryID string) error {
	fileID := uuid.NewRandom().String() + ".jpg"
	filePath := "./source/" + fileID

	fileBuffer, err := getFileBytes(file, filePath)
	if err != nil {
		return fmt.Errorf("unable to get file bytes: %w", err)
	}

	err = s.s3.Upload(bytes.NewReader(fileBuffer), fileID)
	if err != nil {
		return fmt.Errorf("failed to load file to object storage: %w", err)
	}

	err = s.repo.AddFileID(fileID, entryID)
	if err != nil {
		return fmt.Errorf("cannot add file id to the database: %w", err)
	}

	return nil
}

func getFileBytes(file io.ReadCloser, filePath string) ([]byte, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to create file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to copy image to file: %w", err)
	}

	upFile, err := os.Open(filePath)
	if err != nil {
		return []byte{}, fmt.Errorf("could not open local filepath [%v]: %w", filePath, err)
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()

	fileBuffer := make([]byte, upFileInfo.Size())
	upFile.Read(fileBuffer)

	return fileBuffer, nil
}

// GetEntries returns all entries or entries with particular date.
func (s *PhotoService) GetEntries(date string) ([]repository.Entry, error) {
	entries, err := s.repo.GetEntries(date)
	if err != nil {
		return nil, fmt.Errorf("cannot get entries from db: %w", err)
	}

	return entries, nil
}
