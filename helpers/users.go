package helpers

import (
	dbgen "github.com/amanallah-jendoubi/Textio/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
}

func DatabaseUserToUser(dbUser dbgen.User) User {
	var name string
	if dbUser.Name.Valid {
		name = dbUser.Name.String
	}

	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Name:      name,
	}
}
