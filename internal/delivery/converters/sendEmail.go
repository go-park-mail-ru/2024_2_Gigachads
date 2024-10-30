package converters

type SendEmailRequest struct {
	ParentId    int    `json:"parentId"`
	Recipient   string `json:"recipient"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
