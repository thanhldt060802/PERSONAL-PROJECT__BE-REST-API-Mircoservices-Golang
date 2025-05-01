package dto

type SuccessResponse[T any] struct {
	Body struct {
		Status  int    `json:"status" example:"200"`
		Message string `json:"message" example:"Get users successful"`
		Data    T      `json:"data,omitempty"`
		Total   int    `json:"total,omitempty" example:"1"`
	}
}

type ErrorResponse struct {
	Status  int      `json:"status" example:"400"`
	Message string   `json:"message" example:"Get users failed"`
	Error_  string   `json:"error,omitempty" example:"Bad Request"`
	Details []string `json:"details,omitempty" example:"Load from database failed"`
}

func (err *ErrorResponse) Error() string {
	return err.Message
}

func (err *ErrorResponse) GetStatus() int {
	return err.Status
}
