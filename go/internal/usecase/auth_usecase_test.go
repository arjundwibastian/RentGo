package usecase

import (
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"
	"testing"
	"time"
)

func TestRegister_SuccessUser(t *testing.T) {
	mockRepo := &MockUserRepo{}
	fakeAdminKey := "iamAdmin"
	fakeJWTSecret := []byte("super-secret-test-key")
	fakeJWTExpiry := 24 

	uc := NewAuthUseCase(mockRepo,fakeAdminKey,fakeJWTSecret, time.Duration(fakeJWTExpiry))
	req := dto.RegisterRequest{
		Email: "test@mail.com",
		FullName: "test admin",
		Password: "password123",
		PhoneNumber: "081219219",
		Address: "jln setia budi",
		AdminKey: "iamNotAdmin",
	}
	mockRepo.CheckEmailExistsResult = false
	mockRepo.CheckEmailExistsErr = nil


	mockRepo.RegisterUserErr = nil

	result, err := uc.Register(req)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected a user to be returned, got nil")
	}
	if result.Role != "user" {
		t.Fatal("expected user role")
	}
}

func TestRegister_SuccessAdmin(t *testing.T) {
	mockRepo := &MockUserRepo{}
	fakeAdminKey := "iamAdmin"
	fakeJWTSecret := []byte("super-secret-test-key")
	fakeJWTExpiry := 24 

	uc := NewAuthUseCase(mockRepo,fakeAdminKey,fakeJWTSecret, time.Duration(fakeJWTExpiry))
	req := dto.RegisterRequest{
		Email: "test@mail.com",
		FullName: "test admin",
		Password: "password123",
		PhoneNumber: "081219219",
		Address: "jln setia budi",
		AdminKey: "iamAdmin",
	}
	mockRepo.CheckEmailExistsResult = false
	mockRepo.CheckEmailExistsErr = nil
	
	mockRepo.RegisterUserErr = nil

	result, err := uc.Register(req)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected a user to be returned, got nil")
	}
	if result.Role != "admin" {
		t.Fatal("expected admin role")
	}
}


func TestRegister_EmailExists(t *testing.T) {
	mockRepo := &MockUserRepo{}
	fakeAdminKey := "iamAdmin"
	fakeJWTSecret := []byte("super-secret-test-key")
	fakeJWTExpiry := 24 

	uc := NewAuthUseCase(mockRepo,fakeAdminKey,fakeJWTSecret, time.Duration(fakeJWTExpiry))
	req := dto.RegisterRequest{
		Email: "test@mail.com",
		FullName: "test admin",
		Password: "password123",
		PhoneNumber: "081219219",
		Address: "jln setia budi",
		AdminKey: "iamAdmin",
	}
	mockRepo.CheckEmailExistsResult = true
	mockRepo.CheckEmailExistsErr = nil
	
	mockRepo.RegisterUserErr = nil

	result, err := uc.Register(req)
	if err == nil {
		t.Errorf("Expected an error of email taken")
	}
	if err != domain.ErrEmailRegistered {
		t.Errorf("Expected ErrEmailRegistered, got: %v", err)
	}
	if result != nil {
		t.Errorf("Expected user to be nil on failure, got: %v", result)
	}
}

func TestLogin_Success(t *testing.T) {
	mockRepo := &MockUserRepo{}
	fakeAdminKey := "iamAdmin"
	fakeJWTSecret := []byte("super-secret-test-key")
	fakeJWTExpiry := 24 

	uc := NewAuthUseCase(mockRepo,fakeAdminKey,fakeJWTSecret, time.Duration(fakeJWTExpiry))
	req := dto.LoginRequest{
		Email: "test@mail.com",
		Password: "password123",
	}
	mockRepo.ValidateUserLoginResult = &domain.User{ID: 1, Email: "test@mail.com", Role: "admin"}
	mockRepo.ValidateUserLoginErr = nil
	

	result, err := uc.Login(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result == "" {
		t.Errorf("Expected token got nothing")
	}
}

func TestLogin_invalidEmailPassword(t *testing.T) {
	mockRepo := &MockUserRepo{}
	fakeAdminKey := "iamAdmin"
	fakeJWTSecret := []byte("super-secret-test-key")
	fakeJWTExpiry := 24 

	uc := NewAuthUseCase(mockRepo,fakeAdminKey,fakeJWTSecret, time.Duration(fakeJWTExpiry))
	req := dto.LoginRequest{
		Email: "test@mail.com",
		Password: "password123",
	}
	mockRepo.ValidateUserLoginResult = nil
	mockRepo.ValidateUserLoginErr = domain.InvalidEmailorPassword
	

	resultToken, err := uc.Login(req)
	if err == nil {
		t.Errorf("Expected and error but got nothing")
	}
	if err != domain.InvalidEmailorPassword {
		t.Errorf("expecyed invalid email or password erro got %v", err)
	}
	if resultToken != "" {
		t.Errorf("Expected no token")
	}
}

