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
// @Summary Get user profile with BMI
// @Description Retrieves the currently logged-in user profile and BMI calculation using JWT
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response{responseData=dto.BMIReponse}
// @Failure 400 {object} dto.Response "Bad request"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /users/profile [get]
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


