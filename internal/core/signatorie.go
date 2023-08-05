package core

type Sub struct {
	Id     int    `json:"id"`
	Date   string `json:"date" binding:"required"`
	UserId string `json:"userId" db:"user_id"`
	Name   string `json:"name"`
}
