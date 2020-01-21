package pgsql

import (
	"database/sql"

	"srinathkrishna.in/snippetbox/pkg/models"
)

// SnippetModel driver encapsulation
type SnippetModel struct {
	DB *sql.DB
}

// Insert into DB
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	id := 0
	stmt := `INSERT INTO snippets (title, content, created, expires)
			 VALUES ($1, $2, now()::timestamp, now()::timestamp + '3 day'::interval)
			 RETURNING id`
	err := m.DB.QueryRow(stmt, title, content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Get a given snippet
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Latest returns the last few snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
