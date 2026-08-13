package consumer

type CreateConsumer struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	DisplayName string `json:"display_name,omitempty" binding:"omitempty,min=3,max=50"`
	Phone       string `json:"phone,omitempty" binding:"omitempty,e164"`
}
