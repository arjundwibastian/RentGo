package usecase

import (
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
	endDate := req.BookingStart.AddDate(0,0, req.TotalDays - 1)
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
	available, err := uc.repo.GetAvailableVehicle(newBooking.VehicleID, newBooking.BookingStart, newBooking.BookingEnd)
	if err != nil {
		return nil, err
	}
	if (vehicle.Quantity - *available) < 1 {
		return nil, domain.ErrvehicleUnavailable
	}
	updateBalance := user.Balance - totalPrice
	if updateBalance < 0 {
		return nil, domain.ErrNotEnoughBalance
	}
	confirmedBooking, err := uc.repo.CreateBooking(newBooking)
	if err != nil {
		return nil, err
	}
	_ , err = uc.repo.UpdateBalance(user.ID, updateBalance)
	if err != nil {
		return nil, err
	}
	
	err = uc.emailService.SendBookingConfirmation(confirmedBooking, vehicle, user)
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
	booking, err := uc.repo.GetBookingByID(req.BookingID)
	if err != nil {
		return nil, err
	}
	if booking.UserID != userID {
		return nil, domain.ErrInvalidBookingUserID
	}
	if req.Confirm != "confirm" {
		return nil, domain.ErrInvalidConfirmInput
	}
	updatedBooking, err := uc.repo.CancelUserBooking(booking)
	if err != nil {
		return nil, err
	}
	_, err = uc.repo.UpdateBalance(userID, updatedBooking.TotalPrice)
	user, err := uc.repo.GetUserProfile(userID)
	if err != nil {
		return nil, err
	}
	err = uc.emailService.SendBookingCancellation(updatedBooking, user)
	return updatedBooking, nil
}