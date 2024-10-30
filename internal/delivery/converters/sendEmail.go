package converters

type SendEmailRequest struct {
	ID    int    `json:"id"`
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
