package register

import (
	"context"

	"github.com/apotourlyan/ludus-studii/pkg/errorutil"
	"github.com/apotourlyan/ludus-studii/pkg/idutil"
	"github.com/apotourlyan/ludus-studii/pkg/passutil"
	"github.com/apotourlyan/ludus-studii/services/user/internal/domain"
	"github.com/apotourlyan/ludus-studii/services/user/internal/repository"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/errcode"
	"github.com/apotourlyan/ludus-studii/services/user/internal/service/user/register/errtext"
)

type Service struct {
	repo   repository.UserRepository
	idgen  idutil.Generator
	hasher passutil.Hasher
}

func NewService(
	repo repository.UserRepository,
	idgen idutil.Generator,
	hasher passutil.Hasher,
) *Service {
	return &Service{repo, idgen, hasher}
}

func (s *Service) Register(
	ctx context.Context, m *Request,
) (*Response, error) {
	err := Validate(m)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByEmail(ctx, m.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errorutil.NewServiceError(
			errcode.EmailExists,
			errtext.EmailExists)
	}

	passwordHash, err := s.hasher.Hash(m.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:           s.idgen.Next(),
		Email:        m.Email,
		PasswordHash: passwordHash,
		Role:         domain.RoleUser,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	response := &Response{
		ID: user.ID,
	}

	return response, nil
}
