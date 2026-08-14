package usecase

import (
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"
	"time"
)

type MockUserRepo struct {
	CheckEmailExistsResult bool
    CheckEmailExistsErr    error
    RegisterUserErr error
    ValidateUserLoginResult *domain.User
    ValidateUserLoginErr    error
    GetUserProfileResult *domain.User
    GetUserProfileErr    error
    GetUserBalanceResult *int
    GetUserBalanceErr    error
	UpdateBalanceResult	*int
    UpdateBalanceErr error
    GetVehiclesResult *[]domain.Vehicle
    GetVehiclesErr    error
	GetVehiclesAvailableByDateResult *[]dto.VehicleAvailableResponse
	GetVehiclesAvailableByDateErr error
    GetVehicleByIDResult *domain.Vehicle
    GetVehicleByIDErr    error
    GetAvailableVehicleResult *int
    GetAvailableVehicleErr    error
    CreateBookingResult *domain.Booking
    CreateBookingErr    error
	GetUserBookingHistoryResult *[]domain.Booking
	GetUserBookingHistoryErr error
    GetBookingByIDResult *domain.Booking
    GetBookingByIDErr    error
    CheckBookingUserResult bool
    CheckBookingUserErr    error
    CancelUserBookingResult *domain.Booking
    CancelUserBookingErr    error
   

    RegisterUserCalled *domain.User
    AddBalanceCalled   int

	CreateNewVehiclesResult *domain.Vehicle
	CreateNewVehiclesErr error
	UpdateVehiclesResult *domain.Vehicle
	UpdateVehiclesErr error
	GetRevenueResult *dto.RevenueReportResponse
	GetRevenueErr error
	GetTopVehicleResult *[]dto.TopVehicle
	GetTopVehicleErr error

}

func(m *MockUserRepo) CheckEmailExists(email string) (bool,error) {
	return m.CheckEmailExistsResult, m.CheckEmailExistsErr
}
func(m *MockUserRepo) RegisterUser(user *domain.User) error {
	m.RegisterUserCalled = user
	return m.RegisterUserErr
}
func(m *MockUserRepo) ValidateUserLogin(userLogin *dto.LoginRequest) (*domain.User, error) {
	return m.ValidateUserLoginResult, m.ValidateUserLoginErr
}
func(m *MockUserRepo) GetUserProfile(userID int) (*domain.User, error) {
	return m.GetUserProfileResult, m.GetUserProfileErr
}
func(m *MockUserRepo) GetUserBalance(userID int) (*int, error) {
	return m.GetUserBalanceResult, m.GetUserBalanceErr
}
func(m *MockUserRepo) UpdateBalance(userID int, updateBalance int) (*int, error) {
	return m.UpdateBalanceResult, m.UpdateBalanceErr
}
func(m *MockUserRepo) GetVehicles() (*[]domain.Vehicle, error) {
	return m.GetVehiclesResult, m.GetVehiclesErr
}
func(m *MockUserRepo) GetVehiclesAvailableByDate(reqStart, reqEnd time.Time) (*[]dto.VehicleAvailableResponse, error) {
	return m.GetVehiclesAvailableByDateResult, m.GetAvailableVehicleErr
}
func(m *MockUserRepo) GetVehicleByID(vehicleID int) (*domain.Vehicle, error) {
	return m.GetVehicleByIDResult, m.GetVehicleByIDErr
}
func(m *MockUserRepo) GetAvailableVehicle(vehicleID int, start, end time.Time) (*int, error) {
	return m.GetAvailableVehicleResult, m.GetAvailableVehicleErr
}
func(m *MockUserRepo) CreateBooking(booking *domain.Booking) (*domain.Booking, error) {
	return m.CreateBookingResult, m.CreateBookingErr
}
func(m *MockUserRepo) GetUserBookingHistory(userID int) (*[]domain.Booking, error) {
	return m.GetUserBookingHistoryResult, m.GetUserBookingHistoryErr
}
func(m *MockUserRepo) GetBookingByID(bookingID int)(*domain.Booking, error) {
	return m.GetBookingByIDResult, m.GetBookingByIDErr
}
func(m *MockUserRepo) CheckBookingUser(userID, bookingID int)(bool, error) {
	return m.CheckBookingUserResult, m.CheckBookingUserErr
}
func(m *MockUserRepo) CancelUserBooking(booking *domain.Booking) (*domain.Booking, error) {
	return m.CancelUserBookingResult, m.CancelUserBookingErr
}


//admin
func(m *MockUserRepo) CreateNewVehicles(vehicle *domain.Vehicle) (*domain.Vehicle, error) {
	return m.CreateNewVehiclesResult, m.CreateNewVehiclesErr
}
func(m *MockUserRepo) UpdateVehicles(vehicle *domain.Vehicle) (*domain.Vehicle, error) {
	return m.UpdateVehiclesResult, m.UpdateVehiclesErr
}
func(m *MockUserRepo) GetRevenue()(*dto.RevenueReportResponse, error) {
	return m.GetRevenueResult, m.GetRevenueErr
}
func(m *MockUserRepo) GetTopVehicle()(*[]dto.TopVehicle, error) {
	return m.GetTopVehicleResult, m.GetTopVehicleErr
}




type MockEmailService struct {
	SendBookingConfirmationErr error
	SendBookingCancellationErr error
}
func (m *MockEmailService) SendBookingConfirmation(booking *domain.Booking, vehicle *domain.Vehicle, user *domain.User) error {
	return m.SendBookingConfirmationErr
}
func (m *MockEmailService) SendBookingCancellation(booking *domain.Booking, user *domain.User) error {
	return m.SendBookingCancellationErr
}