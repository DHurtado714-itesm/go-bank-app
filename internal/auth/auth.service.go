package auth

import (
	"context"
	"errors"
	"go-bank-app/pkg/jwt"
	"go-bank-app/utils"
	"log"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*User, error)

	// Returns a JWT
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{repo: repo}
}

// Login implements AuthService.
func (s *authService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.HashedPassword) {
		log.Println("❌ Hash error:", err)
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Println("❌ JWT error:", err)
		return "", err
	}

	return token, nil
}

// Register implements AuthService.
func (s *authService) Register(ctx context.Context, email string, password string) (*User, error) {
	existingUser, _ := s.repo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("User already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID: uuid.New().String(),
		Email: email,
		HashedPassword: string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}


