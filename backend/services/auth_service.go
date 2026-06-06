package services

import (
	"Distribyte/backend/database"
	"Distribyte/backend/models"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(
	email string,
	password string,
) (models.User, error) {

	var user models.User

	query := `
	SELECT
		id,
		name,
		email,
		password_hash,
		created_at
	FROM users
	WHERE email = $1
	`

	err := database.DB.QueryRow(
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func RegisterUser(
	name string,
	email string,
	password string,
) error {

	hashedPassword, err :=
		bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)

	if err != nil {
		return err
	}

	query := `
	INSERT INTO users(
		name,
		email,
		password_hash
	)
	VALUES($1,$2,$3)
	`

	_, err = database.DB.Exec(
		query,
		name,
		email,
		string(hashedPassword),
	)

	return err
}
