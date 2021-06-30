package model

import (
	"gofarnay/config"
	"log"
	"time"
)

type User struct {
	ID       	int    `json:"id"`
	Email    	string `json:"email"`
	Password 	string `json:"password"`
	Name     	string `json:"name"`
	CreatedAt	string `json:"create_at"`
	UpdatedAt	string `json:"updated_at"`
	DeletedAt 	string `json:"deleted_at"`
}

func GetUsers(start, count int) ([]User, error) {
	log.Println("Model GetUsers")
	db := config.DbConn()
	defer db.Close()
	rows, err := db.Query("SELECT id, email, password, name FROM users LIMIT ? OFFSET ?", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (u *User) CreateUser() error {
	log.Println("Model CreateUser")
	db := config.DbConn()
	defer db.Close()

	u.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	u.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	_, err := db.Exec("INSERT INTO users(email, password, name, created_at, updated_at) VALUES(?, ?, ?, ?, ?)", u.Email, u.Password, u.Name, u.CreatedAt, u.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUser() error {
	log.Println("Model GetUser")
	db := config.DbConn()
	defer db.Close()

	return db.QueryRow("SELECT email, password, name FROM users where id = ?",
		u.ID).Scan(&u.Email, &u.Password, &u.Name)
}

func (u *User) UpdateUser() error {
	log.Println("Model UpdateUser")
	db := config.DbConn()
	defer db.Close()

	u.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	_, err := db.Exec("UPDATE users SET email=?, password=?, name=?, updated_at=? WHERE id=?",
		u.Email, u.Password, u.Name, u.UpdatedAt, u.ID)

	return err
}

func (u User) DeleteUser() error {
	log.Println("Model DeleteUser")
	db := config.DbConn()
	defer db.Close()

	_, err := db.Exec("DELETE FROM users WHERE id = ? ", u.ID)

	return err
}
