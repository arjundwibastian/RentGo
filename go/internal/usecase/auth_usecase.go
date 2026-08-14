package usecase

import (
	"fmt"
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	userRepo  domain.UserRepository
	adminKey	string
	jwtSecret []byte
	jwtExpiry time.Duration
}

// NewAuthUseCase creates the auth use case with its dependencies injected.
func NewAuthUseCase(userRepo domain.UserRepository, adminKey string, jwtSecret []byte, jwtExpiry time.Duration) domain.AuthUseCase {
	return &authUseCase{
		userRepo:  userRepo,
		adminKey: adminKey,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (uc *authUseCase) Register(req dto.RegisterRequest) (*domain.User, error) {
	exists, err := uc.userRepo.CheckEmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailRegistered
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    req.Email,
		FullName: req.FullName,
		Password: string(hashed),
		Address: req.Address,
		PhoneNumber: req.PhoneNumber,
	}
	if req.AdminKey == uc.adminKey {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}
	if err := uc.userRepo.RegisterUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *authUseCase) Login(req dto.LoginRequest) (string, error) {
	user, err := uc.userRepo.ValidateUserLogin(&req)
	if err != nil {
		return "", err
	}
	token, err := uc.generateToken(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("Error: Failed to Generate Token")
	}

	return token, nil
}

func (uc *authUseCase) generateToken(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role" : role,
		"exp":     time.Now().Add(uc.jwtExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.jwtSecret)
}
