package petitions

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name" binding:"required" db:"name"`
	Grade    string `json:"grade" binding:"required"`
	Password string `json:"password" binding:"required"`
}
