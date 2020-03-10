package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/go-playground/validator.v9"
)

/* validator */

// CustomValidator - validator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate - validator
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

/* others */

// GetTitlesFromCODH - get titles from CODH
func GetTitlesFromCODH() (map[string]string, error) {
	doc, err := goquery.NewDocument("http://codh.rois.ac.jp/char-shape/book/")
	if err != nil {
		return nil, fmt.Errorf("goquery: %v", err)
	}
	log.Printf("INFO: got title list from CODH")

	titles := make(map[string]string)
	doc.Find("#list > tbody > tr").Each(func(i int, s *goquery.Selection) {
		titles[s.AttrOr("id", "")] = s.Find("td:nth-child(3)").Text()
	})

	return titles, err
}
