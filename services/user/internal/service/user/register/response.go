package register

import "github.com/apotourlyan/ludus-studii/services/user/internal/domain"

type Response struct {
	ID   int64
	Role domain.UserRole
}
