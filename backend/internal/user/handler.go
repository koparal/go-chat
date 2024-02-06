package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided data
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body CreateUserData true "User data to register"
// @Success 200 {object} CreateUserResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	var u CreateUserData
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.Register(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Login godoc
// @Summary Login as a user
// @Description Login with the provided credentials
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body LoginUserData true "User login data"
// @Success 200 {object} LoginUserResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var user LoginUserData
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.Service.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", u.AccessToken, 60*60*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, u)
}

// @Summary Logout user
// @Description Logout the user by clearing the JWT cookie
// @Tags users
// @Accept json
// @Produce json
// @Router /logout [post]
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successfully"})
}
