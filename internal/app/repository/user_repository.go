package repository

import (
	"errors"

	"deposit_month/internal/app/ds"
)

var ErrEmailExists = errors.New(
	"пользователь с таким email уже существует",
)

func (r *Repository) CreateUser(
	user *ds.User,
) error {
	var count int64

	err := r.db.
		Model(&ds.User{}).
		Where(
			"email = ?",
			user.Email,
		).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return ErrEmailExists
	}

	return r.db.Create(user).Error
}
