package request

type TranslateClone struct {
	Source      string `json:"source"      validate:"required" example:"auto"`
	Destination string `json:"destination" validate:"required" example:"en"`
	Original    string `json:"original"    validate:"required" example:"Tại vì sao"`
}
