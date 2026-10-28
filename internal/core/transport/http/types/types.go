package types

type Error struct {
	Error string `json:"error"`
}
type Status struct {
	Status string `json:"status" validate:"oneof=Новый Подтверждён Собирается 'Готов к выдаче' Завершён"`
}
