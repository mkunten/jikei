package mojiportal

import (
	"fmt"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

/* routing */

/* struct for mojiportal search query */

// CharSearchParam - char search parameter
type CharSearchParam struct {
	Q       string `query:"q" validate:"required,len=1"`
	Unicode int    `query:"-"`
	Offset  int    `query:"offset" validate:"min=0"`
	Limit   int    `query:"limit" validate:"min=-1"`
	OfEach  bool   `query:"delegate"`
}

/* struct for mojiportal search result */

// List - response json
type List struct {
	SearchResults int    `json:"search_results"`
	List          []Char `json:"list"`
}

// Char - char
type Char struct {
	ID           string `json:"id"`
	Character    string `json:"character"`
	Delegate     bool   `json:"delegate"`
	Unicode      string `json:"unicode"`
	Source       Source `json:"source"`
	ThumbnailURL string `json:"thumbnail_url"`
	ManifestURL  string `json:"manifest_url"`
	Link         string `json:"link"`
	Subject      string `json:"subject"`
	Creator      string `json:"creator"`
	Rights       string `json:"rights"`
	RightsURL    string `json:"rights_url"`
}

// Source - source
type Source struct {
	Label string `json:"label"`
	Title string `json:"title"`
	BID   string `json:"bid"`
	Frame int    `json:"frame"`
	Side  string `json:"side"`
}

// ErrorResponse - error response
type ErrorResponse struct {
	Error []ErrorItem `json:"error"`
}

// ErrorItem - error item
type ErrorItem struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

// GetErrorResponseFromValidateError - get ErrorResponse from validator.v9 error
func GetErrorResponseFromValidateError(e error) ErrorResponse {
	errRes := ErrorResponse{
		Error: []ErrorItem{},
	}
	for _, err := range e.(validator.ValidationErrors) {
		fmt.Printf("%v", err)
		item := ErrorItem{
			Key: strings.ToLower(err.Field()),
			Message: fmt.Sprintf("%s %s but '%v'",
				err.ActualTag(), err.Param(), err.Value()),
		}
		errRes.Error = append(errRes.Error, item)
	}
	return errRes
}
