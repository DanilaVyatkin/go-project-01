package apper

import "encoding/json"

var (
	ErrNotFound = NewAppError(nil, "not found", "", "US-000003")
)

type AppErr struct {
	Err              error  `json:"-"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
	conde            string `json:"conde"`
}

func (e *AppErr) Error() string {
	return e.Message
}

func (e *AppErr) Unwrap() error { return e.Err }

func (e *AppErr) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage, code string) *AppErr {
	return &AppErr{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		conde:            code,
	}
}

func systemError(err error) *AppErr {
	return NewAppError(err, "internal systen err", err.Error(), "US=000000")
}
