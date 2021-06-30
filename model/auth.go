package model

import (
	"github.com/dgrijalva/jwt-go"
	"gofarnay/config"
	"time"
)

type Auth struct {
	Email		string	`json:"email"`
	Password 	string	`json:"password"`
}

type Credentials struct {
	Email		string	`json:"email"`
	Password 	string	`json:"password"`
}

type Claims struct {
	Email	string	`json:"email"`
	jwt.StandardClaims
}

func (c *Credentials) Signin() error {
	db := config.DbConn()
	defer db.Close()

	var u User

	err := db.QueryRow("SELECT id, email, password, name FROM users where email = ? AND password = ?",
		c.Email, c.Password).Scan(&u.ID, &u.Email, &u.Password, &u.Name)

	return err
}

func (c *Credentials) GenerateToken() (string, error) {
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
		Email:  c.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SigningString()
}

//func SetCookie(w http.ResponseWriter, token string) {
//	log.Println("SetCookie")
//	expirationTime := time.Now().Add(time.Minute * 5)
//	cookie := http.Cookie{
//		Name:    "Authorization",
//		Value:   token,
//		Expires: expirationTime,
//	}
//
//	http.SetCookie(w, &cookie)
//}

