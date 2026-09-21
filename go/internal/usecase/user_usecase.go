package usecase

import (
	"log"
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"
)

type userUseCase struct {
	repo    domain.UserRepository
	emailService domain.EmailService
}

func NewUserUseCase(repo domain.UserRepository, emailService domain.EmailService) domain.UserUseCase {
	return &userUseCase{repo: repo, emailService: emailService}
}

func (uc *userUseCase) GetProfile(userID int) (*domain.User, error) {
	user, err := uc.repo.GetUserProfile(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *userUseCase) AddUserTotalBalance(userID int, topUpAmount int) (*int, error) {
	balance, err := uc.repo.GetUserBalance(userID)
	if err != nil {
		return nil, err
	}
	updatedBalance := *balance + topUpAmount
	_,err = uc.repo.UpdateBalance(userID, updatedBalance)
	if err != nil {
		return nil, err
	}
	return &updatedBalance, nil
}

func (uc *userUseCase) GetVehicleList() (*[]domain.Vehicle, error) {
	vehicles, err := uc.repo.GetVehicles()
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}
func (uc *userUseCase) GetVehicleAvailableByDate(req dto.DateRequest) (*[]dto.VehicleAvailableResponse, error) {
	endDate := req.BookingStart.AddDate(0,0, req.TotalDays)
	vehicles, err := uc.repo.GetVehiclesAvailableByDate(req.BookingStart, endDate)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (uc *userUseCase) CreateNewBooking(req dto.BookingRequest, userID int,) (*domain.Booking, error){
	vehicle, err := uc.repo.GetVehicleByID(req.VehicleID)
	if err != nil {
		return nil, err
	}
	user, err := uc.repo.GetUserProfile(userID)
	if err != nil {
		return nil, err
	}
	totalPrice := req.TotalDays * vehicle.DailyRate
	endDate := req.BookingStart.AddDate(0,0, req.TotalDays)
	newBooking := &domain.Booking{
		UserID: userID,
		VehicleID: req.VehicleID,
		BookingStart: req.BookingStart,
		BookingEnd: endDate,
		TotalPrice: totalPrice,
		Status: "confirmed",
	}
	confirmedBooking, err := uc.repo.CreateBookingAtomic(newBooking)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := uc.emailService.SendBookingConfirmation(confirmedBooking, vehicle, user); err != nil {
			log.Printf("send booking confirmation failed booking_id=%d: %v", confirmedBooking.ID, err)
		}
	}()
	return confirmedBooking, nil
}


func (uc *userUseCase) GetBookingHistory(userID int) (*[]domain.Booking, error) {
	bookings, err := uc.repo.GetUserBookingHistory(userID)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (uc *userUseCase) CancelUserBooking(req dto.CancelRequest, userID int) (*domain.Booking, error){
	if req.Confirm != "confirm" {
		return nil, domain.ErrInvalidConfirmInput
	}
	user, err := uc.repo.GetUserProfile(userID)
	if err != nil {
		return nil, err
	}
	updatedBooking, err := uc.repo.CancelBookingAtomic(req.BookingID, userID)
	if err != nil {
		return nil, err
	}

	
	go func() {
		if err := uc.emailService.SendBookingCancellation(updatedBooking, user); err != nil {
			log.Printf("send booking cancellation failed booking_id=%d: %v", updatedBooking.ID, err)
		}
	}()
	return updatedBooking, nil
}