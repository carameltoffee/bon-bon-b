package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strawberry/internal/models"
	"strawberry/internal/service"

	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	FullName       string `json:"full_name"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Specialization string `json:"specialization"`
	ScheduleType   string `json:"schedule_type"`
	Code           string `json:"code"`
}

type RegisterRes struct {
	Id int64 `json:"id"`
}

// @Summary Register a new user or master
// @Description Create a new user or master account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterReq true "registration data"
// @Success 201 {object} RegisterRes
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	var data *RegisterReq
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		newErrorResponse(http.StatusBadRequest, "invalid data", c)
		return
	}
	err := h.s.VerificationCode.VerifyCode(c.Request.Context(), data.Email, data.Code)
	if err != nil {
		newErrorResponse(http.StatusUnauthorized, "invalid code", c)
		return
	}
	if data.Specialization == "" {
		data.Specialization = "user"
	}
	user := &models.User{
		FullName:       data.FullName,
		Email:          data.Email,
		Username:       data.Username,
		Password:       data.Password,
		Specialization: data.Specialization,
		ScheduleType:   data.ScheduleType,
		AverageRating:  5,
	}

	id, err := h.s.Users.Create(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			newErrorResponse(http.StatusInternalServerError, "something on our side", c)
			return
		}
		if errors.Is(err, service.ErrUserExists) {
			newErrorResponse(http.StatusConflict, "user already exists", c)
			return
		}
		var valErr service.ValidationError
		if errors.As(err, &valErr) {
			newErrorResponse(http.StatusBadRequest, valErr.Msg, c)
			return
		}
		newErrorResponse(http.StatusInternalServerError, "unknown error", c)
		return
	}
	c.JSON(http.StatusCreated, &RegisterRes{
		Id: id,
	})
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}

// @Summary Login user
// @Description Authenticate user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginReq true "login credentials"
// @Success 200 {object} LoginRes
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var data *LoginReq
	if err := c.ShouldBindJSON(&data); err != nil {
		newErrorResponse(http.StatusBadRequest, "invalid credentials format", c)
		return
	}

	token, err := h.s.Users.Login(c.Request.Context(), data.Username, data.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			newErrorResponse(http.StatusUnauthorized, "invalid username or password", c)
			return
		}
		newErrorResponse(http.StatusInternalServerError, "auth error", c)
		return
	}

	c.JSON(http.StatusOK, &LoginRes{Token: token})
}

type UpdateReq struct {
	Id             int64  `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Specialization string `json:"specialization"`
	Bio            string `json:"bio"`
	FullName       string `json:"full_name"`
	ScheduleType   string `json:"schedule_type"`
}

// @Summary Update User
// @Description Update User's data such as username , bio,  full_name, specialization and email
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Param input body UpdateReq true "update data"
// @Success 200 "Апдейт"
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	claims, ok := getClaims(c)
	if !ok {
		newErrorResponse(http.StatusUnauthorized, "unauthorized", c)
		return
	}

	var data *UpdateReq
	if err := c.ShouldBindJSON(&data); err != nil {
		newErrorResponse(http.StatusBadRequest, "bad data", c)
		return
	}
	err := h.s.Users.Update(c.Request.Context(), claims.Id, &models.User{
		Id:             data.Id,
		Username:       data.Username,
		FullName:       data.FullName,
		Bio:            data.Bio,
		Specialization: data.Specialization,
		Email:          data.Email,
		ScheduleType:   data.ScheduleType,
		Password:       "pLACEHOLDERPASSWORD_42",
	})
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			newErrorResponse(http.StatusInternalServerError, "something on our side", c)
			return
		}
		if errors.Is(err, service.ErrUserExists) {
			newErrorResponse(http.StatusConflict, "user already exists", c)
			return
		}
		var valErr service.ValidationError
		if errors.As(err, &valErr) {
			newErrorResponse(http.StatusBadRequest, valErr.Msg, c)
			return
		}
	}
	c.Status(http.StatusOK)
}

type RestoreReq struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

// Restore godoc
// @Summary      Restore password
// @Description  Verifies a code and sets a new password for the user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        restore body RestoreReq true "Restore request payload"
// @Success      200 {string} string "OK"
// @Failure      400 {object} ErrorResponse "Invalid data or validation error"
// @Failure      401 {object} ErrorResponse "Invalid code"
// @Failure      404 {object} ErrorResponse "User not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /restore [post]
func (h *Handler) Restore(c *gin.Context) {
	var data *RestoreReq
	if err := c.ShouldBindBodyWithJSON(&data); err != nil {
		newErrorResponse(http.StatusBadRequest, "invalid data", c)
		return
	}
	err := h.s.VerificationCode.VerifyCode(c.Request.Context(), data.Email, data.Code)
	if err != nil {
		newErrorResponse(http.StatusUnauthorized, "invalid code", c)
		return
	}
	if err := h.s.Users.ChangePassword(c.Request.Context(), data.Email, data.Password); err != nil {
		var valErr service.ValidationError
		if errors.As(err, &valErr) {
			newErrorResponse(http.StatusBadRequest, valErr.Msg, c)
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			newErrorResponse(http.StatusNotFound, "that user does not exists", c)
			return
		}
		fmt.Println(err)
		newErrorResponse(http.StatusInternalServerError, "something on our side", c)
		return
	}
	c.Status(http.StatusOK)
}
