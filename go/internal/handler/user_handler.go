package handler

import (
	"errors"
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUC domain.UserUseCase
}

func NewUserHandler(userUC domain.UserUseCase) *UserHandler {
	return &UserHandler{userUC: userUC}
}


// GetProfile godoc
// @Summary Get user profile
// @Description Retrieves the currently logged-in user's profile information using JWT
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.HandlerResponse{responseData=domain.User}
// @Failure 400 {object} dto.HandlerResponse "Bad request"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	user, err := h.userUC.GetProfile(userID)
	if err != nil{
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	return c.JSON(http.StatusOK,dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "",
			ResponseData:    user,
		})
}

// AddUserBalanceHandler godoc
// @Summary Top up user balance
// @Description Adds the specified amount to the currently logged-in user's balance
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.TopUpRequest true "Top Up Data"
// @Success 200 {object} dto.HandlerResponse "Successfully topped up balance"
// @Failure 400 {object} dto.HandlerResponse "Bad request or validation error"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/users/topup [post]
func (h *UserHandler) AddUserBalanceHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	var req dto.TopUpRequest
	if err := c.Bind(&req); err != nil{
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

	currentBalance, err := h.userUC.AddUserTotalBalance(userID, req.Amount)
	if err != nil{
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	return c.JSON(http.StatusOK,dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "successfully top up balance",
			ResponseData:    echo.Map{
				"current balance" : currentBalance,
			},
		})
}

// GetVehicleListHandler godoc
// @Summary Get all vehicles
// @Description Retrieves a list of all available vehicles for rental
// @Tags Vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.HandlerResponse{responseData=[]domain.Vehicle}
// @Failure 400 {object} dto.HandlerResponse "Bad request or validation error"
// @Failure 404 {object} dto.HandlerResponse "No vehicles found"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/vehicles [get]
func (h *UserHandler) GetVehicleListHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	vehicles, err := h.userUC.GetVehicleList()
	if err != nil{
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	if len(*vehicles) == 0 {
		return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: data not found",
				ResponseData:    nil,
			})
	}

	return c.JSON(http.StatusOK,dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "",
			ResponseData:    vehicles,
		})
}

// GetVehicleAvailableByDateHandler godoc
// @Summary Get available vehicles by date
// @Description Retrieves vehicles that are available for booking within the specified date range
// @Tags Vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.DateRequest true "Date Range"
// @Success 200 {object} dto.HandlerResponse{responseData=[]dto.VehicleAvailableResponse}
// @Failure 400 {object} dto.HandlerResponse "Bad request or validation error"
// @Failure 404 {object} dto.HandlerResponse "No vehicles available"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/vehicles/date [post]
func (h *UserHandler) GetVehicleAvailableByDateHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	var req dto.DateRequest
	if err := c.Bind(&req); err != nil{
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

	vehicles, err := h.userUC.GetVehicleAvailableByDate(req)
	if err != nil{
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	if len(*vehicles) == 0 {
		return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: data not found",
				ResponseData:    nil,
			})
	}
	return c.JSON(http.StatusOK,dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "",
			ResponseData:    vehicles,
		})
}

// CreateBookingHandler godoc
// @Summary Create a new booking
// @Description Creates a new vehicle booking for the currently logged-in user
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BookingRequest true "Booking Data"
// @Success 201 {object} dto.HandlerResponse{responseData=domain.Booking}
// @Failure 400 {object} dto.HandlerResponse "Bad request or validation error"
// @Failure 404 {object} dto.HandlerResponse "Vehicle unavailable for selected dates"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/bookings [post]
func (h *UserHandler) CreateBookingHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	var req dto.BookingRequest
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
	newBooking, err := h.userUC.CreateNewBooking(req, userID)
	if err != nil {
		if errors.Is(err, domain.ErrvehicleUnavailable) {
			return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
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
	return c.JSON(http.StatusCreated, dto.HandlerResponse{
		ResponseCode:    "01",
		ResponseMessage: "successfully create new booking",
		ResponseData:    newBooking,
	})
}

// CancelBookingHandler godoc
// @Summary Cancel a booking
// @Description Cancels an existing booking owned by the currently logged-in user
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CancelRequest true "Cancel Data"
// @Success 200 {object} dto.HandlerResponse{responseData=domain.Booking}
// @Failure 400 {object} dto.HandlerResponse "Bad request or invalid confirmation"
// @Failure 401 {object} dto.HandlerResponse "User not authorized to cancel this booking"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/bookings/cancel [post]
func (h *UserHandler) CancelBookingHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	var req dto.CancelRequest
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
	cancelledBooking, err := h.userUC.CancelUserBooking(req, userID)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidBookingUserID) {
			return c.JSON(http.StatusUnauthorized, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		} else if errors.Is(err, domain.ErrInvalidConfirmInput) {
			return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
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
		ResponseCode:    "01",
		ResponseMessage: "successfully cancel booking",
		ResponseData:    cancelledBooking,
	})
}


// GetUserBookingHistoryHandler godoc
// @Summary Get booking history
// @Description Retrieves all bookings made by the currently logged-in user
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.HandlerResponse{responseData=[]domain.Booking}
// @Failure 404 {object} dto.HandlerResponse "No bookings found"
// @Failure 500 {object} dto.HandlerResponse "Internal server error"
// @Router /api/v1/bookings/history [get]
func (h *UserHandler) GetUserBookingHistoryHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)
	if userID == 0 {
		return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "error: user id is invalid",
			ResponseData:    nil,
		})
	}

	bookings, err := h.userUC.GetBookingHistory(userID)
	if err != nil{
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: err.Error(),
			ResponseData:    nil,
		})
	}
	if len(*bookings) == 0 {
		return c.JSON(http.StatusNotFound, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: "Error: data not found",
				ResponseData:    nil,
			})
	}

	return c.JSON(http.StatusOK,dto.HandlerResponse{
			ResponseCode:    "01",
			ResponseMessage: "",
			ResponseData:    bookings,
		})
}


