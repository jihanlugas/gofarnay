package model

import (
	"gofarnay/config"
	"time"
)

type Order struct {
	OrderId   int       `json:"orderId"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	PhotoId   int       `json:"photoId"`
	PhotoUrl  int       `json:"photoUrl"`
	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

func (o *Order) CreateOrder() error {
	db := config.DbConn()
	defer db.Close()

	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	res, err := db.Exec("INSERT INTO orders(name, email, phone, title, message, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?)", o.Name, o.Email, o.Phone, o.Title, o.Message, o.CreatedAt, o.UpdatedAt)

	if err != nil {
		return err
	}

	lid, err := res.LastInsertId()

	if err != nil {
		return err
	}

	o.OrderId = int(lid)

	return nil
}

func (o *Order) UpdatePhotoId() error {
	db := config.DbConn()
	defer db.Close()

	_, err := db.Exec("UPDATE orders set photo_id = ? WHERE order_id = ?", o.PhotoId, o.OrderId)

	return err
}

func GetOrders(start, count int) ([]Order, error) {
	db := config.DbConn()
	defer db.Close()
	rows, err := db.Query("SELECT order_id, name, email, phone, title, message, photo_id, created_at, updated_at FROM orders LIMIT ? OFFSET ?", count, start)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []Order{}

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.OrderId, &o.Name, &o.Email, &o.Phone, &o.Title, &o.Message, &o.PhotoId, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil

}
