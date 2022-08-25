package repository

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
)

// Repository presents type for querying and persisting data.
type Repository struct {
	db *sql.DB
}

// New create a new instance of Repository.
func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// AddFileID adds file id to the database.
func (r *Repository) AddFileID(fileID string) (string, error) {
	var id string

	sql, _, _ := squirrel.
		Insert("files").
		Columns("file_id").
		Values(fileID).
		Suffix("RETURNING id;").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	err := r.db.QueryRow(sql, fileID).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("cannot save file id: %w", err)
	}

	return id, nil
}

// SaveEntry adds particular entry to the database.
func (r *Repository) SaveEntry(copyright, date, explanation, hdURL, mediaType, serviceVersion, title, url, fileID string) (string, error) {
	var id string

	sql, _, _ := squirrel.
		Insert("entries").
		Columns("file_id", "copyright", "date", "explanation", "hd_url", "media_type", "service_version", "title", "url").
		Values(fileID, copyright, date, explanation, hdURL, mediaType, serviceVersion, title, url).
		Suffix("RETURNING id;").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	err := r.db.QueryRow(sql, fileID, copyright, date, explanation, hdURL, mediaType, serviceVersion, title, url).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("cannot save entry: %w", err)
	}

	return id, nil
}

// Entry represents an entry in the form in which it's stored in the database.
type Entry struct {
	ID             string `json:"id"`
	FileID         string `json:"file_id"`
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	HDURL          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
}

// GetEntries retrieves all album entries or entry for the selected day.
func (r *Repository) GetEntries(date string) ([]Entry, error) {
	var entries []Entry

	query := squirrel.
		Select("entries.id, files.file_id, copyright, date, explanation, hd_url, media_type, service_version, title, url, entries.created, entries.updated").
		From("entries").
		Join("files ON entries.file_id = files.id").
		PlaceholderFormat(squirrel.Dollar)

	if date != "" {
		query = query.Where(squirrel.Eq{"date": date})
	}

	sql, args, _ := query.ToSql()

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error with query executing: %w", err)
	}

	for rows.Next() {
		entry := Entry{}

		if err := rows.Scan(
			&entry.ID,
			&entry.FileID,
			&entry.Copyright,
			&entry.Date,
			&entry.Explanation,
			&entry.HDURL,
			&entry.MediaType,
			&entry.ServiceVersion,
			&entry.Title,
			&entry.URL,
			&entry.Created,
			&entry.Updated,
		); err != nil {
			return nil, fmt.Errorf("cannot get entries: %w", err)
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with result set: %w", err)
	}

	return entries, nil
}
