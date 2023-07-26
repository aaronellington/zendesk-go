package zendesk

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	StatusCode  int    `json:"status_code"`
	Body        []byte `json:"body"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (err *Error) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("Zendesk API Error, Status Code: %d", err.StatusCode)
	}

	return err.Message
}

func (err *Error) UnmarshalJSON(b []byte) error {
	err1 := errorResponse1{}
	if err1err := json.Unmarshal(b, &err1); err1err == nil {
		if err1.Error.Message != "" {
			err.Message = err1.Error.Title
			err.Description = err1.Error.Message

			return nil
		}
	}

	err2 := errorResponse2{}
	if err2err := json.Unmarshal(b, &err2); err2err == nil {
		if err2.Error != "0" {
			err.Message = err2.Error
			err.Description = err2.Description
		}

		return nil
	}

	return nil
}

type errorResponse1 struct {
	Error error1 `json:"error"`
}

type error1 struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type errorResponse2 struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}
