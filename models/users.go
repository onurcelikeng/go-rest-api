
package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	u "go-contacts/utils"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token";sql:"-"`
}

func (user *User) Validate() (map[string] interface{}, bool) {
	if !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(user.Password) < 8 {
		return u.Message(false, "Password is required"), false
	}

	tempUser := &User{}

	err := GetDB().Table("users").Where("email = ?", user.Email).First(tempUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if tempUser.Email != "" {
		return u.Message(false, "Email already in use."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (user *User) Create() (map[string] interface{}) {
	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	GetDB().Create(user)
	if user.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{UserId: user.ID})
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString
	user.Password = ""

	response := u.Message(true, "Account has been created")
	response["user"] = user
	return response
}

func Login(email, password string) (map[string]interface{}) {
	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &Token{UserId: user.ID})
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString
	user.Password = ""

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

func GetUser(u uint) *User {
	user := &User{}
	GetDB().Table("users").Where("id = ?", u).First(user)
	if user.Email == "" {
		return nil
	}

	user.Password = ""
	return user
}
