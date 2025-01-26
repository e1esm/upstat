package repository

import (
	"database/sql"
	"fmt"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
	"github.com/chamanbravo/upstat/svcerr"
)

func (r *Repository) SaveUser(u *dto.UserSignUp) error {
	stmt, err := r.db.Prepare("INSERT INTO users(username, email, password) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Username, u.Email, u.Password)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

func (r *Repository) FindUserByUsernameAndEmail(u *dto.UserSignUp) (*models.User, error) {
	stmt, err := r.db.Prepare("SELECT id, username, email FROM users WHERE username = $1 OR email = $2")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user := new(models.User)
	result := stmt.QueryRow(u.Username, u.Email).Scan(&user.ID, &user.Username, &user.Email)
	if result != nil {
		if result == sql.ErrNoRows {
			return nil, svcerr.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to scan rows for users: %w", err)
	}
	return user, nil
}

func (r *Repository) FindUserByUsernameAndPassword(username, password string) (*models.User, error) {
	stmt, err := r.db.Prepare("SELECT id, username, email, firstname, lastname FROM users WHERE username = $1 AND password = $2")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user := new(models.User)
	err = stmt.QueryRow(username, password).Scan(&user.ID, &user.Username, &user.Email, &user.Firstname, &user.Lastname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, svcerr.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to find user by username and password: %w", err)
	}
	return user, nil
}

func (r *Repository) FindUserByUsername(username string) (*models.User, error) {
	stmt, err := r.db.Prepare("SELECT id, username, email, firstname, lastname, password FROM users WHERE username = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user := new(models.User)
	result := stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Email, &user.Firstname, &user.Lastname, &user.Password)
	if result != nil {
		if result == sql.ErrNoRows {
			return nil, svcerr.ErrUserNotFound
		}

		return nil, err
	}
	return user, nil
}

func (r *Repository) UsersCount() (int, error) {
	stmt, err := r.db.Prepare("SELECT COUNT(*) FROM users")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("failed to get users count: %w", err)
	}

	return count, nil
}

func (r *Repository) UpdatePassword(username string, newPassword string) error {
	stmt, err := r.db.Prepare("UPDATE users SET password = $1 WHERE username = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newPassword, username)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (r *Repository) UpdateAccount(username string, u *dto.UpdateAccountIn) error {
	stmt, err := r.db.Prepare("UPDATE users SET firstname = $1, lastname = $2 WHERE username = $3")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Firstname, u.Lastname, username)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}
