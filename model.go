package main

import (
	"fmt"
	"strings"

	"github.com/labstack/gommon/log"
)

// Jikei - csv items of dataset
type Jikei struct {
	JID     string `db:"jid,primarykey" json:"jid"`
	Unicode int    `db:"unicode,notnull" json:"unicode"`
	PID     string `db:"pid,notnull" json:"pid"`
	X       int    `db:"x,notnull" json:"x"`
	Y       int    `db:"y,notnull" json:"y"`
	BlockID string `db:"block_id,notnull" json:"block_id"`
	CharID  string `db:"char_id,notnull" json:"char_id"`
	Width   int    `db:"width,notnull" json:"width"`
	Height  int    `db:"height,notnull" json:"height"`
	Created string `db:"created,notnull" json:"created"`
}

// Pages - pages
type Pages struct {
	List []Page
}

// Page - csv items of pages
type Page struct {
	PID     string `db:"pid,primarykey" json:"pid"`
	BID     string `db:"bid,notnull" json:"bid"`
	Title   string `db:"title,notnull" json:"title"`
	Pos     int    `db:"pos,notnull" json:"pos"`
	Frame   int    `db:"frame,notnull" json:"frame"`
	Side    string `db:"side,notnull" json:"side"`
	Width   int    `db:"width,notnull" json:"width"`
	Height  int    `db:"height,notnull" json:"height"`
	Created string `db:"created,notnull" json:"created"`
}

// GetCanvasID - get IIIF canvas @id
func (page *Page) GetCanvasID(iiif *IIIF, id string) string {
	return fmt.Sprintf("%s/%s/%s/canvas/%s",
		iiif.PAPIBase, page.BID, page.PID, id)
}

// GetCanvasIDWithChars - get IIIF canvas @id with queried chars
func (page *Page) GetCanvasIDWithChars(iiif *IIIF, id, chars string) string {
	return fmt.Sprintf("%s/%s/%s/%s/canvas/%s",
		iiif.PAPIBase, page.BID, page.PID, url.PathEscape(chars), id)
}

// GetImageID - get IIIF image @id
func (page *Page) GetImageID(iiif *IIIF) string {
	return fmt.Sprintf("%s/%s/images/%s.jpg", iiif.IAPIBase, page.BID, page.PID)
}

// GetLabel - get label
func (page *Page) GetLabel() string {
	return fmt.Sprintf("%s %d%s", target.Title, target.Frame, target.Side)
}

// FormatPageManifest - format page manifest
func (page *Page) FormatPageManifest(config *Config, chars string) iiif.Manifest {
	base := config.baseURL + "/api/page/" + page.PID
	otherContents := []iiif.OtherContent{}
	if chars != "" {
		base += "/" + url.PathEscape(chars)
		otherContents = append(otherContents, iiif.OtherContent{
			ID:   base + "/canvas/c1/annolist",
			Type: "sc:AnnotationList",
		})
	}
	jpg := page.GetIIIFImageID()

	manifest := iiif.Manifest{
		Context: "http://iiif.io/api/presentation/2/context.json",
		ID:      base + "/manifest",
		Type:    "sc:Manifest",
		Label: fmt.Sprintf("『%s』（%s）%d%s",
			page.Title, page.BID, page.Frame, page.Side),
		Metadata: []iiif.Metadatum{
			iiif.Metadatum{
				Label: "BID",
				Value: page.BID,
			},
			iiif.Metadatum{
				Label: "TITLE",
				Value: page.Title,
			},
			iiif.Metadatum{
				Label: "FRAME",
				Value: strconv.Itoa(page.Frame),
			},
			iiif.Metadatum{
				Label: "SIDE",
				Value: page.Side,
			},
		},
		ViewingHint: "paged",
		License:     manifestLicense,
		Attribution: manifestAttribution,
		Logo:        baseURL + "/img/nijl_symbolmark.jpg",
		Related: []iiif.IDFormat{
			iiif.IDFormat{
				ID: fmt.Sprintf("https://kotenseki.nijl.ac.jp/biblio/%s/viewer/%d",
					page.BID, page.Frame),
				Format: "text/html",
			},
		},
		Within: base + "/",
		Sequences: []iiif.Sequence{
			iiif.Sequence{
				ID:   base + "/sequence",
				Type: "sc:Sequence",
				Canvases: []iiif.Canvas{
					iiif.Canvas{
						ID:     base + "/canvas/c1",
						Type:   "sc:Canvas",
						Label:  fmt.Sprintf("%d%s", page.Frame, page.Side),
						Width:  page.Width,
						Height: page.Height,
						Images: []iiif.Image{
							iiif.Image{
								ID:         base + "/canvas/c1/annotion/anno1",
								Type:       "oa:Annotation",
								Motivation: "sc:painting",
								Resource: iiif.Resource{
									ID:     jpg + "/full/full/0/default.jpg",
									Type:   "dctypes:Image",
									Format: "image/jpeg",
									Width:  page.Width,
									Height: page.Height,
									Service: iiif.Service{
										Context: "http://iiif.io/api/image/2/context.json",
										ID:      jpg,
										Profile: "http://iiif.io/api/image/2/level1.json",
									},
								},
								On: base + "/canvas/c1",
							},
						},
						OtherContent: otherContents,
					},
				},
			},
		},
	}
	return manifest
}

// JikeiPages - jikei page views
type JikeiPages struct {
	List []JikeiPageView
}

// JikeiPageView - jikei page view
type JikeiPageView struct {
	JID      string `db:"jid,primarykey" json:"jid"`
	Unicode  int    `db:"unicode,notnull" json:"unicode"`
	PID      string `db:"pid,notnull" json:"pid"`
	X        int    `db:"x,notnull" json:"x"`
	Y        int    `db:"y,notnull" json:"y"`
	BlockID  string `db:"block_id,notnull" json:"block_id"`
	CharID   string `db:"char_id,notnull" json:"char_id"`
	Width    int    `db:"width,notnull" json:"width"`
	Height   int    `db:"height,notnull" json:"height"`
	Created  string `db:"created,notnull" json:"created"`
	BID      string `db:"bid,notnull" json:"bid"`
	Title    string `db:"title,notnull" json:"title"`
	Pos      int    `db:"pos,notnull" json:"pos"`
	Frame    int    `db:"frame,notnull" json:"frame"`
	Side     string `db:"side,notnull" json:"side"`
	PWidth   int    `db:"pwidth,notnull" json:"pwidth"`
	PHeight  int    `db:"pheight,notnull" json:"pheight"`
	PCreated string `db:"pcreated,notnull" json:"pcreated"`
}

const jpvKeys = "jid, unicode, j.pid, x, y, block_id, char_id, j.width, j.height, j.created, bid, title, pos, frame, side, p.width AS pwidth, p.height AS pheight, p.created AS pcreated"

// GetIIIFImageID - get IIIF image @id from id
func (jp *JikeiPageView) GetIIIFImageID() string {
	return fmt.Sprintf("%s/%s/characters/U+%X/%s.jpg",
		iiifBaseURL, jp.BID, jp.Unicode, jp.JID)
}

// ResponseListUpload - response json format
type ResponseListUpload struct {
	Message  string `json:"message"`
	Total    int    `json:"total"`
	Inserted int    `json:"inserted"`
	Updated  int    `json:"updated"`
	Error    int    `json:"error"`
}

// QuerySearchParam - query search parameters
type QuerySearchParam struct {
	Q      string `query:"q" validate:"required,min=1"`
	Offset int    `query:"offset" validate:"min=0"`
	Limit  int    `query:"limit" validate:"min=-1"`
}

// QuerySearchParamQ - parsed query of query search parameters
type QuerySearchParamQ struct {
	Chars []rune
	BID   []string
}

// ParseQuery - QuerySearchParam.ParseQuery
func (param *QuerySearchParam) ParseQuery() (qspq QuerySearchParamQ) {
	cols := strings.Split(param.Q, " ")
	for _, col := range cols {
		if pos := strings.Index(col, ":"); pos != -1 {
			key := col[0:pos]
			value := col[pos+1:]
			switch key {
			case "chars":
				qspq.Chars = append(qspq.Chars, []rune(value)...)
			case "bid":
				qspq.BID = append(qspq.BID, value)
			default:
				log.Errorf("QuerySearchParam.ParseQuery: unexpected key: %s", key)
			}
		} else {
			qspq.Chars = append(qspq.Chars, []rune(col)...)
		}
	}
	return qspq
}

// QuerySearchList - query search list
type QuerySearchList struct {
	Total  int                   `json:"total"`
	Offset int                   `json:"offset"`
	Limit  int                   `json:"limit"`
	Prev   string                `json:"prev,omitempty"`
	Next   string                `json:"next,omitempty"`
	List   []QuerySearchListChar `json:"list"`
}

// QuerySearchListChar - char
type QuerySearchListChar struct {
	ID           string                `json:"id"`
	Character    string                `json:"character"`
	Unicode      string                `json:"unicode"`
	Source       QuerySearchListSource `json:"source"`
	ThumbnailURL string                `json:"thumbnail_url"`
	ManifestURL  string                `json:"manifest_url"`
	Link         string                `json:"link"`
}

// QuerySearchListSource - source
type QuerySearchListSource struct {
	Label  string `json:"label"`
	Title  string `json:"title"`
	BID    string `json:"bid"`
	Pos    int    `json:"pos"`
	Frame  int    `json:"frame"`
	Side   string `json:"side"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
