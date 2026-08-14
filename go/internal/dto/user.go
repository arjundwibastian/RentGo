package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullName" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	PhoneNumber	string `json:"phoneNumber" validate:"required,numeric,min=8,max=15"`
	Address  string `json:"address" validate:"required"`
	AdminKey string `json:"adminKey"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TopUpRequest struct {
	Amount int `json:"amount" validate:"required,min=10000"`
}

type DateRequest struct {
	BookingStart time.Time `json:"bookingStart" validate:"required"`
	TotalDays	int	`json:"totalDays" validate:"required,min=1"`
}
type BookingRequest struct {
	VehicleID    	int `json:"vehicleID" validate:"required"`
	TotalDays		int	`json:"totalDays" validate:"required,min=1"`
	BookingStart 	time.Time `json:"bookingStart" validate:"required"`
}
type CancelRequest struct {
	BookingID  int `json:"bookingID" validate:"required"`
	Confirm		string `json:"confirm" validate:"required,eq=confirm"`
}
type VehicleRequest struct {
	Name	string `json:"name" validate:"required"`
	Description	string `json:"description" validate:"required"`
	Quantity	int `json:"quantity" validate:"required"`
	DailyRate	int `json:"dailyRate" validate:"required"`
	Category	string `json:"category" validate:"required"`
}