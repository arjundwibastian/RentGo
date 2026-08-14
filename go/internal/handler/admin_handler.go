package handler

import (
	"errors"
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminUC domain.AdminUseCase
}

func NewAdminHandler(adminUC domain.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUC: adminUC}
}


// CreateNewVehicleHandler godoc
// @Summary Create a new vehicle
// @Description Creates a new vehicle to the vehicles database
// @Tags Vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.VehicleRequest true "Vehicle Data"
// @Success 201 {object} dto.HandlerResponse{responseData=domain.Vehicle}
// @Failure 400 {object} dto.HandlerResponse "Bad request or validation error"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/Vehicles [post]
func (h *AdminHandler) CreateNewVehicleHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	var req dto.VehicleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
		})
	}
	newVehicle, err := h.adminUC.CreateNewVehicle(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	return c.JSON(http.StatusCreated, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "successfully create new Vehicle",
		ResponseData:    newVehicle,
	})
}


// UpdateVehicleHandler godoc
// @Summary Update existing vehicle
// @Description update existing vehicle database
// @Tags Vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "vehicle ID"
// @Param request body dto.VehicleRequest true "Vehicle Data"
// @Success 200 {object} dto.HandlerResponse{responseData=domain.Vehicle}
// @Failure 400 {object} dto.HandlerResponse "Bad request or validation error"
// @Failure 404 {object} dto.HandlerResponse "vehicle data is not found"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/Vehicles/{id} [put]
func (h *AdminHandler) UpdateVehicleHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	vehicleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode: "01",
			ResponseMessage: err.Error(),
			ResponseData: nil,
		})
		
	}

	var req dto.VehicleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
		})
	}
	updatedVehicle, err := h.adminUC.UpdateVehiclesById(req, vehicleID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "ERROR: vehicles data not found",
			ResponseData:    nil,
		})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
		}

	}
	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "successfully update vehicle",
		ResponseData:    updatedVehicle,
	})
}


// GetRevenueReportHandler godoc
// @Summary Get revenue report
// @Description retrieves revenu report for all completed booking from database
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.HandlerResponse{responseData=dto.RevenueReportResponse}
// @Failure 400 {object} dto.HandlerResponse "Bad request or validation error"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/reports/revenue [GET]
func (h *AdminHandler) GetRevenueReportHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	
	report, err := h.adminUC.GetRevenueReport()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
		ResponseCode:    "01",
		ResponseMessage: err.Error(),
		ResponseData:    nil,
		})
		

	}
	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "",
		ResponseData:    report,
	})
}

// GetTopVehicleHandler godoc
// @Summary Get top vehicle report
// @Description retrieves report for top most booked vehicle from bookings database
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.HandlerResponse{responseData=dto.TopVehicle}
// @Failure 400 {object} dto.HandlerResponse "Bad request or validation error"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/reports/top-vehicle [GET]
func (h *AdminHandler) GetTopVehicleHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	
	report, err := h.adminUC.GetTopVehicle()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
		ResponseCode:    "01",
		ResponseMessage: err.Error(),
		ResponseData:    nil,
		})
		

	}
	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "00",
		ResponseMessage: "",
		ResponseData:    report,
	})
}
