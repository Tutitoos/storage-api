package domain

type IResultError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type IResultData[T any] struct {
	Message   string         `json:"message,omitempty"`
	Data      *T             `json:"data,omitempty"`
	Errors    []IResultError `json:"errors"`
	HasErrors bool           `json:"hasErrors"`
}

func ResultData[T any]() *IResultData[T] {
	return &IResultData[T]{
		Errors:    []IResultError{},
		HasErrors: false,
	}
}

func (r *IResultData[T]) AddMessage(message string) {
	r.Message = message
}

func (r *IResultData[T]) AddData(data T) {
	r.Data = &data
}

func (r *IResultData[T]) AddError(code int, message string) {
	r.Errors = append(r.Errors, IResultError{
		Code:    code,
		Message: message,
	})
	r.HasErrors = true
}
