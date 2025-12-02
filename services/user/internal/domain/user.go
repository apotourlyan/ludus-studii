package domain

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         UserRole
}
