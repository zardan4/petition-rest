package petitions

type User struct {
	Id       int    `json:"id" swaggerignore:"true"`
	Name     string `json:"name" binding:"required" db:"name"`
	Grade    string `json:"grade" binding:"required"`
	Password string `json:"password" binding:"required"`
}
