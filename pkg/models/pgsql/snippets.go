package pgsql

import (
	"database/sql"
	"fmt"
	"strconv"

	"srinathkrishna.in/snippetbox/pkg/models"
)

// SnippetModel driver encapsulation
type SnippetModel struct {
	DB *sql.DB
}

// Insert into DB
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	id := 0

	// this is required because QueryRow does not interpolate the interval string since
	// it is quoted
	expiry, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		return 0, err
	}

	stmt := fmt.Sprintf(`INSERT INTO snippets (title, content, created, expires)
			 VALUES ($1, $2, now()::timestamp, now()::timestamp + '%v day'::interval)
			 RETURNING id`, expiry)
	err = m.DB.QueryRow(stmt, title, content).Scan(&id)
	if err != nil {
		return 0, err
	} else if err == sql.ErrNoRows {
		return 0, models.ErrNoRecord
	}

	return id, nil
}

// Get a given snippet
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}

	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > now()::timestamp AND id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	}

	return s, nil
}

// Latest returns the last few snippets
func (m *SnippetModel) Latest(limit int) ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > now()::timestamp ORDER BY created DESC LIMIT $1`

	rows, err := m.DB.Query(stmt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
