package domain

import (
	"milestone-02/internal/dto"
	"time"
)

type UserRepository interface {
	CheckEmailExists(email string) (bool, error)
	RegisterUser(user *User) error
	ValidateUserLogin(userLogin *dto.LoginRequest) (*User, error)
	GetUserProfile(userID int) (*User, error)
	GetUserBalance(userID int) (*int64, error)
	UpdateBalance(userID int, updateBalance int64) (*int64, error)
	AddUserBalanceAtomic(userID int, amount int64) (*int64, error)
	GetVehicles() (*[]Vehicle, error)
	GetVehiclesAvailableByDate(reqStart, reqEnd time.Time) (*[]dto.VehicleAvailableResponse, error)
	GetVehicleByID(vehicleID int) (*Vehicle, error)
	GetAvailableVehicle(vehicleID int, start, end time.Time) (*int, error)
	CreateBooking(booking *Booking) (*Booking, error)
	CreateBookingAtomic(booking *Booking) (*Booking, error)
	GetUserBookingHistory(userID int) (*[]Booking, error)
	GetBookingByID(bookingID int)(*Booking, error) 
	CheckBookingUser(userID, bookingID int)(bool, error)
	CancelUserBooking(booking *Booking) (*Booking, error)
	CancelBookingAtomic(bookingID, userID int) (*Booking, error)
	CompletePastBookings() error
	


	//admin repo
	CreateNewVehicles(vehicle *Vehicle) (*Vehicle, error)
	UpdateVehicles(vehicle *Vehicle) (*Vehicle, error)
	GetRevenue()(*dto.RevenueReportResponse, error)
	GetTopVehicle()(*[]dto.TopVehicle, error)
}


type EmailService interface {
    SendBookingConfirmation(booking *Booking, vehicle *Vehicle, user *User) error
    SendBookingCancellation(booking *Booking, user *User) error 
}