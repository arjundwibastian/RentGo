package usecase

import (
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"
)

type adminUseCase struct {
	repo domain.UserRepository
}

func NewAdminUsecase(repo domain.UserRepository) domain.AdminUseCase {
	return &adminUseCase{repo: repo}
}

func (ac *adminUseCase) CreateNewVehicle(req dto.VehicleRequest)(*domain.Vehicle, error) {
	vehicle := &domain.Vehicle{
		Name: req.Name,
		Description: req.Description,
		Quantity: req.Quantity,
		DailyRate: req.DailyRate,
		Category: req.Category,
	}
	newVehicle, err := ac.repo.CreateNewVehicles(vehicle)
	if err != nil {
		return nil, err
	}
	return newVehicle, nil
}

func (ac *adminUseCase) UpdateVehiclesById(req dto.VehicleRequest, vehicleID int) (*domain.Vehicle, error) {
	vehicle := &domain.Vehicle{
		ID: vehicleID,
		Name: req.Name,
		Description: req.Description,
		Quantity: req.Quantity,
		DailyRate: req.DailyRate,
		Category: req.Category,
	}
	updateVehicle, err := ac.repo.UpdateVehicles(vehicle)
	if err != nil {
		return nil, err
	}
	return updateVehicle, nil
}

func(ac *adminUseCase) GetRevenueReport() (*dto.RevenueReportResponse, error) {
	report, err := ac.repo.GetRevenue()
	if err != nil {
		return nil, err
	}
	return report, nil
}
func(ac *adminUseCase) GetTopVehicle() (*[]dto.TopVehicle, error) {
	topReport, err := ac.repo.GetTopVehicle()
	if err != nil {
		return nil, err
	}
	return topReport, nil
}