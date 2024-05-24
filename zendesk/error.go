package zendesk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Error struct {
	Response    *http.Response `json:"response"`
	Message     string         `json:"message"`
	Description string         `json:"description"`
}

func (err *Error) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("Zendesk API Error, Status Code: %d", err.Response.StatusCode)
	}

	return err.Message
}

func (err *Error) ImmutableRecord() bool {
	if strings.HasPrefix(err.Message, "RecordInvalid") ||
		strings.HasPrefix(err.Message, "RecordNotFound") {
		return true
	}

	return false
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
		if err2.Error != "" {
			err.Message = err2.Error
			err.Description = err2.Description

			if len(err2.Details) > 0 {
				details := []string{}

				for errorKey, errorDetails := range err2.Details {
					for _, errorDetail := range errorDetails {
						details = append(
							details,
							fmt.Sprintf(
								"[%s: %s - %s]",
								errorKey,
								errorDetail.Error,
								errorDetail.Description,
							),
						)
					}
				}

				err.Message = fmt.Sprintf(
					"%s. Error details: %s",
					err.Message,
					strings.Join(details, ", "),
				)
			}
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
	Error       string                   `json:"error"`
	Description string                   `json:"description"`
	Details     map[string][]ErrorDetail `json:"details"`
}

type ErrorDetail struct {
	Description string `json:"description"`
	Error       string `json:"error"`
}
