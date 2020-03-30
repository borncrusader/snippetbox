package pgsql

import (
	"database/sql"
	"fmt"

	"srinathkrishna.in/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, saltedPassword, salt string) error {
	var id string

	stmt := fmt.Sprintf(`INSERT INTO users (name, email, salted_password, salt, created)
						 VALUES ($1, $2, $3, $4, now()::timestamp) RETURNING ID`)
	err := m.DB.QueryRow(stmt, name, email, saltedPassword, salt).Scan(&id)
	if err != nil {
		return err
	} else if err == sql.ErrNoRows {
		return models.ErrNoRecord
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	return "", nil
}

func (m *UserModel) Get(id string) (*models.User, error) {
	return nil, nil
}
