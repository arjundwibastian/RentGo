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
