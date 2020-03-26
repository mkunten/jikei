package iiif

/* struct */

// IIIF - iiif struct with settings
type IIIF struct {
	IAPIBase string
	PAPIBase string
	RootDir  string
}

// ManifestOptions - IIIF Manifest constructor options
type ManifestOptions struct {
	Prefix string
	Chars  string
}

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
	Structures       []Structure `json:"structures,omitempty"`
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

// Structure - IIIF Manifest Structure
type Structure struct {
	ID          string   `json:"@id"`
	Type        string   `json:"@type"`
	Label       string   `json:"label"`
	ViewingHint string   `json:"viewingHint,omitempty"`
	Members     []Member `json:"members"`
}

// Member - IIIF Manifest Member
type Member struct {
	ID    string `json:"@id"`
	Type  string `json:"@type"`
	Label string `json:"label"`
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

/* interface */

// Target - target
type Target interface {
	GetID() string
	GetIIIFImageID() string
}

/* constructor */

// NewIIIF - IIIF constructor
func NewIIIF(i, p, r string) IIIF {
	return IIIF{
		IAPIBase: i,
		PAPIBase: p,
		RootDir:  r,
	}
}

/* method */

// GetManifest - get manifest file from targets
func (iiif *IIIF) GetManifest(targets []Target, opts *ManifestOptions) {
	canvases := make([]iiif.Canvas, 0, len(targets))
	chars := opts.Chars
	for _, target := range targets {
		id := target.GetID()
		var canvas string
		var otherContents []iiif.OtherContent
		if chars == "" {
			canvas = target.GetCanvasID(iiif, "c1")
		} else {
			canvas = target.GetCanvasIDWithChars(iiif, "c1", chars)
			otherContents = append(otherContents, iiif.OtherContent{
				ID:   canvas + "/annolist",
				Type: "sc:AnnotationList",
			})
		}

		jpg := target.GetIIIFImageID()

		canvases = append(canvases, iiif.Canvas{
			ID:     canvas,
			Type:   "sc:Canvas",
			Label:  target.GetLabel(),
			Width:  target.Width,
			Height: target.Height,
			Images: []iiif.Image{
				iiif.Image{
					ID:         canvas + "/annotion/anno1",
					Type:       "oa:Annotation",
					Motivation: "sc:painting",
					Resource: iiif.Resource{
						ID:     jpg + "/full/full/0/default.jpg",
						Type:   "dctypes:Image",
						Format: "image/jpeg",
						Width:  target.Width,
						Height: target.Height,
						Service: iiif.Service{
							Context: "http://iiif.io/api/image/2/context.json",
							ID:      jpg,
							Profile: "http://iiif.io/api/image/2/level1.json",
						},
					},
					On: canvas,
				},
			},
			OtherContent: otherContents,
		})
	}

	manifest := iiif.Manifest{
		Context: "http://iiif.io/api/presentation/2/context.json",
		ID: fmt.Sprintf("%s/api/manifest?q=%s",
			iiif.IAPIBase, url.QueryEscape(chars)),
		Type:             "sc:Manifest",
		Label:            chars,
		ViewingHint:      "paged",
		ViewingDirection: "right-to-left",
		License:          opts.License,
		Attribution:      opts.Attribution,
		Logo:             iiif.IAPIBase + "/img/nijl_symbolmark.jpg",
		Related: []iiif.IDFormat{
			iiif.IDFormat{
				ID:     "https://kotenseki.nijl.ac.jp/",
				Format: "text/html",
			},
		},
		Within: iiif.IAPIBase + "/",
		Sequences: []iiif.Sequence{
			iiif.Sequence{
				ID: fmt.Sprintf("%s/api/sequence?q=%s",
					iiif.IAPIBase, url.QueryEscape(chars)),
				Type:     "sc:Sequence",
				Canvases: canvases,
			},
		},
	}

	return manifest
}
