package domain

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name" validate:"required,max=100"`
	Spec     string    `json:"spec" validate:"required,max=300"`
	Category string    `json:"category" validate:"oneof=Смартфоны Ноутбуки Аудио Часы"`
	Price    int64     `json:"price" validate:"gt=0,lte=10000000"`
	Old      int64     `json:"old" validate:"gte=0,lte=10000000"`
	Stock    int       `json:"stock" validate:"gte=0,lte=100000"`
	Type     string    `json:"type" validate:"oneof=phone laptop headphones watch"`
	Color    string    `json:"color" validate:"required,hexcolor"`
	Badge    string    `json:"badge" validate:"max=40"`
}
type OrderItem struct {
	ID    uuid.UUID `json:"id" validate:"required"`
	Qty   int       `json:"qty" validate:"gte=1,lte=100"`
	Name  string    `json:"name"`
	Price int64     `json:"price"`
}
type OrderInput struct {
	Name     string      `json:"name" validate:"required,max=80"`
	Phone    string      `json:"phone" validate:"required,min=10,max=20"`
	Delivery string      `json:"delivery" validate:"oneof=pickup delivery"`
	Address  string      `json:"address" validate:"required_if=Delivery delivery,max=200"`
	Items    []OrderItem `json:"items" validate:"required,min=1,max=100,dive"`
}
type Order struct {
	ID uuid.UUID `json:"id"`
	OrderInput
	Total     int64     `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
