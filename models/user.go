package models

import (
	"crypto/bcrypt"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"html"
	"os"
	"strings"
	"time"

	"myshipper/utils/token"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"size:255;not null;unique" json:"username"`
	Password  string `gorm:"size:255;not null;" json:"password"`
	Roles     []Role `gorm:"many2many:users_roles;"`
	FirstName string `gorm:"varchar(255);not null"`
	LastName  string `gorm:"varchar(255);not null"`
	Email     string `gorm:"column:email;unique_index"`
}

func (user *User) SaveUser() (*User, error) {
	var err error
	err = DB.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, error) {
	var err error
	u := User{}
	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error
	if err != nil {
		return "", err
	}
	err = VerifyPassword(password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	token, err := token.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetUserByID(uid uint) (User, error) {
	var u User
	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}
	u.PrepareGive()
	return u, nil
}

func (user *User) PrepareGive() {
	user.Password = ""
}

func (user *User) IsAdmin() bool {
	for _, role := range user.Roles {
		if role.Name == "ROLE_ADMIN" {
			return true
		}
	}
	return false
}

func (user *User) IsNotAdmin() bool {
	return !user.IsAdmin()
}

func (user *User) GenerateJwtToken() string {
	jwtToken := jwt.New(jwt.SigningMethodHS512)
	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}
	jwtToken.Claims = jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 24 * 90).Unix(),
	}
	returnToken, _ := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return returnToken
}

func (user *User) IsValidPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(user.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}
