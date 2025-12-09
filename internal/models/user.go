package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)


type User struct {
	ID uuid.UUID 
	Username string 
	Email string 
	Password string 
	Role string 
	Name string 
	CreatedAt time.Time 
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}