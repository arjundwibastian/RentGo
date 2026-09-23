package domain

import (
	"milestone-02/internal/dto"
)

type AuthUseCase interface {
	Register(req dto.RegisterRequest) (*User, error)
	Login(req dto.LoginRequest) (string, error)
}

type UserUseCase interface {
	GetProfile(userID int) (*User, error)
	AddUserTotalBalance(userID int, topUpAmount int64) (*int64, error)
	GetVehicleList() (*[]Vehicle, error)
	GetVehicleAvailableByDate(req dto.DateRequest) (*[]dto.VehicleAvailableResponse, error)
	CreateNewBooking(req dto.BookingRequest, userID int) (*Booking, error)
	GetBookingHistory(userID int) (*[]Booking, error)
	CancelUserBooking(req dto.CancelRequest, userID int) (*Booking, error)
	// GetWorkoutsByID(userID int, workoutID int) (*Workout, error)
	// CreateNewWorkouts(req dto.WorkoutRequest, userID int) (*Workout, error)
	// UpdateWorkoutsByID(userID int, workoutID int, req dto.WorkoutRequest) (*Workout, error)
	// DeleteWorkoutsByID(userID int, workoutID int) error
	// CreateNewExercise(req dto.ExerciseRequest, userID int) (*Exercise, error)
	// DeleteExerciseByID(userID int, exerciseID int) error
	// CreateNewExerciseLogs(userID int, req dto.ExerciseLogsRequest) (*ExerciseLog, error)
	// GetUserExerciseLogs(userID int) (*[]ExerciseLog, error)
}

type AdminUseCase interface {
	CreateNewVehicle(req dto.VehicleRequest)(*Vehicle, error)
	UpdateVehiclesById(req dto.VehicleRequest, vehicleID int) (*Vehicle, error)
	//report
	GetRevenueReport() (*dto.RevenueReportResponse, error)
	GetTopVehicle()	(*[]dto.TopVehicle, error)
}