package usecase

import (
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestAddUserBalance_Success(t *testing.T) {
	mockRepo := &MockUserRepo{}
	mockEmail := &MockEmailService{}

	uc := NewUserUseCase(mockRepo, mockEmail)
	req := dto.TopUpRequest{
		Amount: 200000,
	}
	balance := int64(150000)
	userID := 1
	mockRepo.GetUserBalanceResult = &balance
	mockRepo.GetUserBalanceErr = nil

	updateBalance := req.Amount + balance
	mockRepo.UpdateBalanceResult = &updateBalance
	mockRepo.UpdateBalanceErr = nil
	result, err := uc.AddUserTotalBalance(userID, req.Amount)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected a result to be returned, got nil")
	}
	if *result != updateBalance {
		t.Fatal("expected correct balance got")
	}
}

func TestAddUserBalance_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepo{}
	mockEmail := &MockEmailService{}

	uc := NewUserUseCase(mockRepo, mockEmail)
	req := dto.TopUpRequest{
		Amount: 200000,
	}
	balance := int64(150000)
	userID := 7377
	mockRepo.GetUserBalanceResult = nil
	mockRepo.GetUserBalanceErr = gorm.ErrRecordNotFound

	updateBalance := req.Amount + balance
	mockRepo.UpdateBalanceResult = &updateBalance
	mockRepo.UpdateBalanceErr = nil
	result, err := uc.AddUserTotalBalance(userID, req.Amount)
	if err == nil {
		t.Errorf("Expected error")
	}
	if err != gorm.ErrRecordNotFound {
		t.Fatal("Expected record not found error")
	}
	if result != nil {
		t.Fatal("expected no result")
	}
}

func TestCreateNewBooking_Success(t *testing.T) {
	mockRepo := &MockUserRepo{}
	mockEmail := &MockEmailService{}

	uc := NewUserUseCase(mockRepo, mockEmail)

	req:= dto.BookingRequest {
		VehicleID: 1,
		TotalDays: 2,
		BookingStart: time.Now(),
	}
	userID := 1

	mockRepo.GetVehicleByIDResult =&domain.Vehicle{ID: 1, DailyRate: 100000, Quantity: 2}
	mockRepo.GetAvailableVehicleErr = nil

	mockRepo.GetUserProfileResult = &domain.User{ID: 1, Balance: 500000}
	mockRepo.GetUserProfileErr = nil
	
	currentlyBooked := 0
	mockRepo.GetAvailableVehicleResult = &currentlyBooked
	mockRepo.GetAvailableVehicleErr = nil

	expectedBooking := &domain.Booking{ID: 10, TotalPrice: 200000, Status: "confirmed"}
	mockRepo.CreateBookingResult = expectedBooking
	mockRepo.CreateBookingErr = nil
	mockRepo.CreateBookingAtomicResult = expectedBooking
	mockRepo.CreateBookingAtomicErr = nil

	updatedBalance := int64(300000)
	mockRepo.UpdateBalanceResult = &updatedBalance
	mockRepo.UpdateBalanceErr = nil 

	mockEmail.SendBookingConfirmationErr = nil
	
	result, err := uc.CreateNewBooking(req, userID)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result == nil {
		t.Fatal("Expected a result to be returned, got nil")
	}
	if result.TotalPrice != 200000 {
		t.Fatal("expected correct price")
	}

}

func TestCreateNewBooking_NotEnoughBalance(t *testing.T) {
	mockRepo := &MockUserRepo{}
	mockEmail := &MockEmailService{}

	uc := NewUserUseCase(mockRepo, mockEmail)

	req:= dto.BookingRequest {
		VehicleID: 1,
		TotalDays: 3,
		BookingStart: time.Now(),
	}
	userID := 1

	mockRepo.GetVehicleByIDResult =&domain.Vehicle{ID: 1, DailyRate: 200000, Quantity: 2}
	mockRepo.GetAvailableVehicleErr = nil

	mockRepo.GetUserProfileResult = &domain.User{ID: 1, Balance: 500000}
	mockRepo.GetUserProfileErr = nil
	
	currentlyBooked := 0
	mockRepo.GetAvailableVehicleResult = &currentlyBooked
	mockRepo.GetAvailableVehicleErr = nil
	mockRepo.CreateBookingAtomicResult = nil
	mockRepo.CreateBookingAtomicErr = domain.ErrNotEnoughBalance

	result, err := uc.CreateNewBooking(req, userID)

	if err == nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if err != domain.ErrNotEnoughBalance {
		t.Fatal("Expected a result to be returned, got nil")
	}
	if result != nil{
		t.Fatal("expected correct price")
	}

}
