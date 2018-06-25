package db

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"gitlab.com/bytesized/bytesized-streaming/helpers"
	"io/ioutil"
	"math/rand"
	"path"
	"time"
)

type CommonModelFields struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type User struct {
	UUIDable
	CommonModelFields
	Login        string `gorm:"not null;unique" json:"login"`
	Admin        bool   `gorm:"not null" json:"admin"`
	PasswordHash string `gorm:"not null" json:"-"`
	Salt         string `gorm:"not null" json:"-"`
}

type Invite struct {
	gorm.Model
	Code   string
	UserID uint
	User   User
}

func (self *User) ValidPassword(password string) bool {
	ctx.Db.Where("login = ?", self.Login).Find(self)
	if self.HashPassword(password, self.Salt) == self.PasswordHash {
		return true
	} else {
		return false
	}
}

func (self *User) SetPassword(password string, salt string) string {
	self.Salt = salt
	self.PasswordHash = self.HashPassword(password, self.Salt)

	return self.PasswordHash
}

func (self *User) HashPassword(password string, salt string) string {
	h := sha256.New()
	h.Write([]byte(salt))
	h.Write([]byte(password))
	hashedStr := hex.EncodeToString(h.Sum(nil))
	return hashedStr
}

// TODO Maran: Create a way to return all errors at once
func CreateUser(login string, password string, admin bool, code string) (User, error) {
	invite := Invite{}
	if len(login) < 3 {
		return User{}, fmt.Errorf("Login should be at least 3 characters")
	}

	if len(password) < 8 {
		return User{}, fmt.Errorf("Password should be at least 8 characters")
	}

	if code != "" {
		count := 0
		ctx.Db.Where("code = ? and user_id IS NULL", code).Find(&invite).Count(&count)
		if count != 0 {
			fmt.Println("Valid and unused code, creating account")
		} else {
			fmt.Println("Not a valid code or already used")
			return User{}, fmt.Errorf("Invite code invalid")
		}
	}

	user := User{Login: login, Admin: admin}
	user.SetPassword(password, randString(24))
	dbobj := ctx.Db.Create(&user)
	if !ctx.Db.NewRecord(&user) {
		invite.UserID = user.ID
		ctx.Db.Save(&invite)
	}
	return user, dbobj.Error
}

func AllUsers() (users []User) {
	ctx.Db.Find(&users)
	return users
}

type UserClaims struct {
	Login  string `json:"login"`
	UserID uint   `json:"user_id"`
	Admin  bool   `json:"admin"`
	jwt.StandardClaims
}

// TODO Maran: Rotate secrets
func TokenSecret() (string, error) {
	tokenPath := path.Join(helpers.GetHome(), ".config", "bss", "token.secret")
	if helpers.FileExists(tokenPath) {
		secret, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			return "", err
		} else {
			return string(secret), nil
		}
	} else {
		secret := randString(32)
		ioutil.WriteFile(tokenPath, []byte(secret), 0700)
		return secret, nil
	}
}

// TODO Maran: Consider setting the jti if we want to increase security.
func (self *User) CreateJWT() (string, error) {
	expiresAt := time.Now().Add(time.Hour * 24).Unix()

	claims := UserClaims{self.Login, self.ID, self.Admin, jwt.StandardClaims{ExpiresAt: expiresAt, Issuer: "bss"}}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := TokenSecret()
	if err != nil {
		return "", err
	}

	ss, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return ss, nil
}

// Plucked from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func randString(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}
