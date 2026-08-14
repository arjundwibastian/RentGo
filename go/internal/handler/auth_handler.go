package handler

import (
	"errors"
	"milestone-02/internal/domain"
	"milestone-02/internal/dto"

	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUC domain.AuthUseCase
}

func NewAuthHandler(authUC domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// RegisterHandler godoc
// @Summary Register a new User
// @Description Creates a new user account in the database
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration Data"
// @Success 200 {object} dto.Response{responseData=domain.User}
// @Failure 400 {object} dto.Response "Bad Request"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /users/register [post]
func (h *AuthHandler) RegisterHandler(c echo.Context) error {
	var req dto.RegisterRequest
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
	user, err := h.authUC.Register(req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailRegistered) {
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
		ResponseMessage: "successfully register user",
		ResponseData:    user,
	})
}

// LoginHandler godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} dto.Response{responseData=object}
// @Failure 400 {object} dto.Response "Bad request"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /users/login [post]
func (h *AuthHandler) LoginHandler(c echo.Context) error {
	var req dto.LoginRequest
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

	token, err := h.authUC.Login(req)
	if err != nil {
		if errors.Is(err, domain.InvalidEmailorPassword) {
			return c.JSON(http.StatusBadRequest, dto.HandlerResponse{
				ResponseCode:    "01",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.HandlerResponse{
				ResponseCode:    "00",
				ResponseMessage: err.Error(),
				ResponseData:    nil,
			})
		}
	}
	return c.JSON(http.StatusOK, dto.HandlerResponse{
		ResponseCode:    "01",
		ResponseMessage: "successfully login",
		ResponseData: echo.Map{
			"token" : token,
		},
	})

}
