package service

type CreateService struct {
	Name     string  `json:"name" binding:"required,min=3,max=30"`
	Desc     string  `json:"desc" binding:"required,min=10,max=100"`
	Duration int32   `json:"duration" binding:"required,gt=0,lte=1440"`
	Charge   float64 `json:"charge" binding:"required,gt=0"`
}

type UpdateService struct {
	Name     string  `json:"name" binding:"omitempty,min=3,max=30"`
	Desc     string  `json:"desc" binding:"omitempty,min=10,max=100"`
	Duration int32   `json:"duration" binding:"omitempty,gt=0,lte=1440"`
	Charge   float64 `json:"charge" binding:"omitempty,gt=0"`
}
