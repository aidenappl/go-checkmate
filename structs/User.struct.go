package structs

import "time"

type UserRole string

const (
	UserRoleOwner  UserRole = "owner"
	UserRoleAdmin  UserRole = "admin"
	UserRoleMember UserRole = "member"
)

type User struct {
	ID         int       `db:"id" json:"id"`
	Email      string    `db:"email" json:"email"`
	Name       *string   `db:"name" json:"name,omitempty"`
	Role       UserRole  `db:"role" json:"role"` // owner, admin, member
	InsertedAt time.Time `db:"inserted_at" json:"inserted_at"`
	Password   string    `db:"password" json:"-"`
}
