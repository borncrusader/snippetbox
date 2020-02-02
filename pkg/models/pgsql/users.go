package pgsql

import (
	"database/sql"

	"srinathkrishna.in/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, salted_password, salt string) error {
	return nil
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	return "", nil
}

func (m *UserModel) Get(id string) (*models.User, error) {
	return nil, nil
}
