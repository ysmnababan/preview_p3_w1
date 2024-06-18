package handler

import (
	"fmt"
	"log"
	"net/http"
	"preview/logger"
	"preview/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func (r *Repo) Register(c echo.Context) error {
	var req models.User
	err := c.Bind(&req)
	if err != nil {
		logger.Logging(c).Error(err)
		return c.JSON(http.StatusInternalServerError, "error binding")
	}

	// validate inputted data
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return c.JSON(http.StatusBadRequest, "invalid param")
	}

	hashedpwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	req.Password = string(hashedpwd)
	res := r.DB.Create(&req)
	if res.Error != nil {
		logger.Logging(c).Error(res.Error)
		return c.JSON(http.StatusInternalServerError, "error creating data")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"username": req.Username, "email": req.Email})
}

func generateToken(s models.User) (string, error) {
	payload := jwt.MapClaims{
		"email":   s.Email,
		"user_id": s.UserID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte("secret-key"))
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("unable to get token")
	}
	return tokenString, nil
}

func (r *Repo) Login(c echo.Context) error {
	var req models.User
	err := c.Bind(&req)
	if err != nil {
		logger.Logging(c).Error(err)
		return c.JSON(http.StatusInternalServerError, "error binding")
	}

	// validate inputted data
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, "invalid param")
	}
	var u models.User

	res := r.DB.Where("email = ?", req.Email).First(&u)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusBadRequest, "error pwd or email")
		}
		logger.Logging(c).Error(res.Error)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	// check pwd
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "error pwd or email")
	}

	token, err := generateToken(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")

	}

	return c.JSON(http.StatusOK, map[string]interface{}{"token": token})
}

func (r *Repo) Loan(c echo.Context) error {
	user_id := uint(c.Get("user_id").(float64))

	var u models.User
	res := r.DB.First(&u, user_id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusBadRequest, "no user found")
		}
		logger.Logging(c).Error(res.Error)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	// create new limit
	u.Limit = 1000.0
	u.Balance = 1000.0
	res = r.DB.Save(&u)
	if res.Error != nil {
		logger.Logging(c).Error(res.Error)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"limit": u.Limit})
}

func (r *Repo) Limit(c echo.Context) error {
	user_id := uint(c.Get("user_id").(float64))

	var u models.User
	r.DB.First(&u, user_id)
	return c.JSON(http.StatusOK, map[string]interface{}{"Limit": u.Limit, "Balance": u.Balance})
}

func (r *Repo) DrawBalance(c echo.Context) error {
	user_id := uint(c.Get("user_id").(float64))

	var u models.User
	r.DB.First(&u, user_id)

	var req models.User
	err := c.Bind(&req)
	if err != nil {
		logger.Logging(c).Error(err)
		return c.JSON(http.StatusInternalServerError, "error binding")
	}

	// validate inputted data
	if req.Balance <= 0 {
		return c.JSON(http.StatusBadRequest, "invalid param")
	}

	if req.Balance > u.Balance {
		return c.JSON(http.StatusBadRequest, "insufficient balance")
	}

	u.Balance = u.Balance - req.Balance
	r.DB.Save(&u)

	return c.JSON(http.StatusOK, map[string]interface{}{"Remaining balance": u.Balance})
}

func (r *Repo) Pay(c echo.Context) error {
	user_id := uint(c.Get("user_id").(float64))

	var u models.User
	r.DB.First(&u, user_id)

	var req models.User
	err := c.Bind(&req)
	if err != nil {
		logger.Logging(c).Error(err)
		return c.JSON(http.StatusInternalServerError, "error binding")
	}

	// validate inputted data
	if req.Balance <= 0 {
		return c.JSON(http.StatusBadRequest, "invalid param")
	}

	if req.Balance+u.Balance > u.Limit {
		return c.JSON(http.StatusBadRequest, "limit exceeded")
	}

	u.Balance = u.Balance + req.Balance
	r.DB.Save(&u)

	return c.JSON(http.StatusOK, map[string]interface{}{"Updated balance": u.Balance})
}
