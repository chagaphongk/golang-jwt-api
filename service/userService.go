package service

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/chagaphongk/register-api/constant"
	"github.com/chagaphongk/register-api/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

type Handler struct {
	DB *mgo.Session
}

func (h *Handler) Signup(c echo.Context) (err error) {
	// Bind
	u := &model.User{
		ID:       bson.NewObjectId(),
		CreateAt: time.Now(),
	}
	if err = c.Bind(u); err != nil {
		return
	}

	// Validate
	user, r := h.getUser(u.Email)
	if r != nil {
		return r
	}

	if user.Email != "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Your emial was used, Please try with new account"}
	}

	if u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or password"}
	}

	encPassword := sha256.Sum256([]byte(u.Password))
	u.Password = fmt.Sprintf("%x", encPassword)

	// Save user
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(constant.DBName).C(constant.DBCollection).Insert(u); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, u)
}

func (h *Handler) Login(c echo.Context) (err error) {
	// Bind
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return
	}

	encPassword := sha256.Sum256([]byte(u.Password))
	u.Password = fmt.Sprintf("%x", encPassword)

	fmt.Println(u.Password)

	// Find user
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB(constant.DBName).C(constant.DBCollection).
		Find(bson.M{"email": u.Email, "password": u.Password}).One(u); err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid email or password"}
		}
		return
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	u.Token, err = token.SignedString([]byte(constant.Key))
	if err != nil {
		return err
	}

	u.Password = "" // Don't send password
	return c.JSON(http.StatusOK, u)
}

func (h *Handler) getUser(email string) (model.User, error) {
	var user model.User

	// Find user
	db := h.DB.Clone()
	defer db.Close()
	if err := db.DB(constant.DBName).C(constant.DBCollection).
		Find(bson.M{"email": email}).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return user, nil
		}
		return user, err
	}

	return user, nil
}
