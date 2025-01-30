package app

import (
	"fmt"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
)

func (a *App) FindUserByUsernameAndEmail(user *dto.UserSignUp) (*models.User, error) {
	u, err := a.db.FindUserByUsernameAndEmail(user)
	if err != nil {
		return nil, fmt.Errorf("faield to find user by email: %w", err)
	}

	return u, nil
}

func (a *App) SaveUser(u *dto.UserSignUp) error {
	return a.db.SaveUser(u)
}

func (a *App) FindUserByUsername(username string) (*models.User, error) {
	u, err := a.db.FindUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by its username: %w", err)
	}

	return u, nil
}

func (a *App) UsersCount() (int, error) {
	return a.db.UsersCount()
}

func (a *App) UpdatePassword(username string, newPassword string) error {
	return a.db.UpdatePassword(username, newPassword)
}

func (a *App) UpdateAccount(username string, u *dto.UpdateAccountIn) error {
	return a.db.UpdateAccount(username, u)
}
