package iiif

/* struct for IIIF manifest files */

// Manifest - IIIF Manifest
type Manifest struct {
	Context          string      `json:"@context"`
	ID               string      `json:"@id"`
	Type             string      `json:"@type,omitempty"`
	Label            string      `json:"label"`
	Metadata         []Metadatum `json:"metadata,omitempty"`
	Description      string      `json:"description,omitempty"`
	ViewingHint      string      `json:"viewingHint,omitempty"`
	ViewingDirection string      `json:"viewingDirection,omitempty"`
	License          string      `json:"license,omitempty"`
	Attribution      string      `json:"attribution,omitempty"`
	Logo             string      `json:"logo,omitempty"`
	SeeAlso          []IDFormat  `json:"seeAlso,omitempty"`
	Related          []IDFormat  `json:"related,omitempty"`
	Within           string      `json:"within,omitempty"`
	Sequences        []Sequence  `json:"sequences,omitempty"`
}

// Metadatum - IIIF Manifest Metadata
type Metadatum struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Sequence - IIIF Manifest Sequence
type Sequence struct {
	ID       string   `json:"@id"`
	Type     string   `json:"@type"`
	Canvases []Canvas `json:"canvases"`
}

// Canvas - IIIF Manifest Canvas
type Canvas struct {
	ID           string         `json:"@id"`
	Type         string         `json:"@type"`
	Label        string         `json:"label"`
	Width        int            `json:"width"`
	Height       int            `json:"height"`
	Images       []Image        `json:"images"`
	OtherContent []OtherContent `json:"otherContent,omitempty"`
}

// Image - IIIF Manifest Image
type Image struct {
	ID         string   `json:"@id"`
	Type       string   `json:"@type"`
	Motivation string   `json:"motivation"`
	Resource   Resource `json:"resource"`
	On         string   `json:"on"`
}

// Resource - IIIF Manifest Resource
type Resource struct {
	ID      string  `json:"@id"`
	Type    string  `json:"@type"`
	Format  string  `json:"format"`
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	Service Service `json:"service"`
}

// Service - IIIF Manifest Service
type Service struct {
	Context string `json:"@context"`
	ID      string `json:"@id"`
	Profile string `json:"profile"`
}

// IDFormat - IIIF Manifest IDFormat
type IDFormat struct {
	ID     string `json:"@id"`
	Format string `json:"format"`
}

// OtherContent - IIIF Manifest OtherContent
type OtherContent struct {
	ID   string `json:"@id"`
	Type string `json:"type"`
}

/* struct for annotation list */

// AnnoList - anno list
type AnnoList struct {
	ID        string          `json:"@id"`
	Context   string          `json:"@context"`
	Type      string          `json:"@type"`
	Resources []AnnoResources `json:"resources"`
}

// AnnoResources - anno resource
type AnnoResources struct {
	ID         string         `json:"@id"`
	Context    string         `json:"@context"`
	Type       string         `json:"@type"`
	Motivation []string       `json:"motivation"`
	Resource   []AnnoResource `json:"resource"`
	On         string         `json:"on"`
	// On         []AnnoOn       `json:"on"`
}

// AnnoResource - anno resource resource
type AnnoResource struct {
	Type   string `json:"@type"`
	Format string `json:"format,omitempty"`
	Chars  string `json:"chars"`
}

// AnnoOn -anno on
// type AnnoOn struct {
// 	Type     string       `json:"@type"`
// 	Full     string       `json:"full"`
// 	Selector AnnoSelector `json:"selector"`
// }
//
// // AnnoSelector - anno selector
// type AnnoSelector struct {
// 	Type    string      `json:"@type"`
// 	Default AnnoDefault `json:"default"`
// 	Item    AnnoItem    `json:"item,omitempty"`
// }
//
// // AnnoDefault - anno default
// type AnnoDefault struct {
// 	Type  string `json:"@type"`
// 	Value string `json:"value"`
// }
//
// // AnnoItem - anno item
// type AnnoItem struct {
// 	Type  string `json:"@type"`
// 	Value string `json:"value"`
// }
