package petitions

type Petition struct {
	Id      int    `json:"id"`
	Title   string `json:"title" binding:"required"`
	Date    string `json:"date" binding:"required"`
	Timeend string `json:"timeend" binding:"required"`
	Text    string `json:"text" binding:"required"`
	Answer  string `json:"answer"`
}

type UpdatePetitionInput struct {
	Title   *string `json:"title"`
	Date    *string `json:"date"`
	Timeend *string `json:"timeend"`
	Text    *string `json:"text"`
	Answer  *string `json:"answer"`
}
