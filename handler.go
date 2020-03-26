package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/mkunten/jikei/pkg/iiif"
	"github.com/mkunten/jikei/pkg/mojiportal"
)

// QueryManifestParam - query manifest parameters
type QueryManifestParam struct {
	Q      string `query:"q" validate:"required"`
	Offset int    `query:"offset" validate:"min=0"`
	Limit  int    `query:"limit" validate:"min=-1"`
	Paged  bool   `query:"paged"`
	Page   int    `query:"page" validate:"min=0"`
}

// GetQueryManifest - get manifest by queries
func GetQueryManifest(c echo.Context) error {
	param := new(QueryManifestParam)
	if err := c.Bind(param); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	runes := []rune(param.Q)
	q := ""
	for _, r := range runes {
		q += "," + strconv.Itoa(int(r))
	}
	q = q[1:]

	var (
		pages []Page
		sql   string
	)

	if param.Paged {
		if param.Page > 0 {
			param.Page--
		}
		param.Limit = 100
		param.Offset = param.Page * param.Limit
		sql = "SELECT distinct(page.*) FROM page INNER JOIN jikei ON page.pid = jikei.pid AND jikei.unicode IN (" + q + ") ORDER BY page.pid OFFSET :Offset LIMIT :Limit"
	} else if param.Limit > 0 {
		sql = "SELECT distinct(page.*) FROM page INNER JOIN jikei ON page.pid = jikei.pid AND jikei.unicode IN (" + q + ") ORDER BY page.pid OFFSET :Offset LIMIT :Limit"
	} else {
		sql = "SELECT distinct(page.*) FROM page INNER JOIN jikei ON page.pid = jikei.pid AND jikei.unicode IN (" + q + ") ORDER BY page.pid"
	}

	_, err := dbmap.Select(&pages, sql, param)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, "resource not found: "+param.Q)
	}

	return c.JSON(http.StatusOK, FormatManifest(pages, param))
}

// GetBiblioManifest - get biblio manifest by bid
func GetBiblioManifest(c echo.Context) error {
	bid := c.Param("bid")
	chars := c.Param("chars")

	var pages []Page
	_, err := dbmap.Select(&pages,
		"SELECT * FROM page WHERE bid = $1 ORDER BY pid", bid)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, "resource not found: "+bid)
	}

	return c.JSON(http.StatusOK, FormatBiblioManifest(pages, chars))
}

// GetPageManifest - get page manifest by page id
func GetPageManifest(c echo.Context) error {
	pid := c.Param("pid")
	chars := c.Param("chars")

	var page Page
	err := dbmap.SelectOne(&page, "SELECT * FROM page WHERE pid = $1", pid)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, "resource not found: "+pid)
	}

	return c.JSON(http.StatusOK, page.FormatPageManifest(chars))
}

// GetPageAnnoList - get page annolist by image
func GetPageAnnoList(c echo.Context) error {
	pid := c.Param("pid")
	chars := c.Param("chars")

	var jikeis []JikeiPageView
	if chars == "" {
		_, err := dbmap.Select(&jikeis,
			"SELECT * FROM jikei WHERE pid = $1 ORDER BY jid", pid)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, "resource not found: "+pid)
		}
	} else {
		runes := []rune(chars)
		q := ""
		for _, r := range runes {
			q += "," + strconv.Itoa(int(r))
		}
		q = q[1:]
		sql := "SELECT * FROM jikei WHERE pid = $1 AND unicode IN (" +
			q + ") ORDER BY jid"
		_, err := dbmap.Select(&jikeis, sql, pid)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusNotFound, "resource not found: "+pid+":"+chars)
		}
	}

	return c.JSON(http.StatusOK, FormatPageAnnoList(jikeis, pid, chars))
}

// GetCharManifest - get char manifest by char id
func GetCharManifest(c echo.Context) error {
	jid := c.Param("jid")

	var j JikeiPageView
	err := dbmap.SelectOne(&j, "SELECT "+jpvKeys+" FROM jikei AS j INNER JOIN page AS p ON j.pid = p.pid WHERE jid = $1", jid)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, "resource not found: "+jid)
	}

	return c.JSON(http.StatusOK, FormatCharManifest(j))
}

// GetQuerySearch - get search by queries
func GetQuerySearch(c echo.Context) error {
	param := new(QuerySearchParam)
	if err := c.Bind(param); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := c.Validate(param); err != nil {
		return c.JSON(http.StatusBadRequest,
			mojiportal.GetErrorResponseFromValidateError(err))
	}

	params := param.ParseQuery()

	q := ""
	for _, r := range []rune(params.Chars) {
		q += "," + strconv.Itoa(int(r))
	}
	q = q[1:]

	if param.Limit == 0 {
		param.Limit = 20
	}

	// count
	sql := "SELECT count(*) FROM jikei AS j INNER JOIN page AS p ON j.pid = p.pid WHERE unicode IN (" + q + ")"
	total, err := dbmap.SelectInt(sql)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// data
	var jikeis []JikeiPageView

	sql = "SELECT " + jpvKeys + " FROM jikei AS j INNER JOIN page AS p ON j.pid = p.pid WHERE unicode IN (" + q + ") ORDER BY jid"
	if param.Offset != 0 {
		sql += " OFFSET :Offset"
	}
	if param.Limit != -1 {
		sql += " LIMIT :Limit"
	}
	_, err = dbmap.Select(&jikeis, sql, param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	chars := make([]QuerySearchListChar, 0, len(jikeis))

	for _, j := range jikeis {
		ch := string(j.Unicode)
		chars = append(chars, QuerySearchListChar{
			ID:        j.JID,
			Character: ch,
			Unicode:   fmt.Sprintf("%04X", j.Unicode),
			Source: QuerySearchListSource{
				Label: fmt.Sprintf("『%s』（%s）%d%s",
					j.Title, j.BID, j.Frame, j.Side),
				BID:    j.BID,
				Title:  j.Title,
				Pos:    j.Pos,
				Frame:  j.Frame,
				Side:   j.Side,
				X:      j.X,
				Y:      j.Y,
				Width:  j.Width,
				Height: j.Height,
			},
			ThumbnailURL: j.GetIIIFImageID() + "/full/full/0/default.jpg",
			ManifestURL:  baseURL + "/api/char/" + j.JID + "/manifest",
			Link: fmt.Sprintf("%s/viewer/?manifest=%s/api/biblio/%s/%s/manifest&pos=%d&xywh=%d,%d,%d,%d",
				baseURL, baseURL, j.BID, url.PathEscape(ch),
				j.Pos, j.X, j.Y, j.Width, j.Height),
		})
	}

	var prev, next string
	if param.Offset >= param.Limit {
		offset := param.Offset - param.Limit
		if offset < 0 {
			offset = 0
		}
		prev = fmt.Sprintf("%s/api/search?q=%s&offset=%d&limit=%d",
			baseURL, url.QueryEscape(param.Q), offset, param.Limit)
	}
	if param.Offset+param.Limit < int(total) {
		next = fmt.Sprintf("%s/api/search?q=%s&offset=%d&limit=%d",
			baseURL, url.QueryEscape(param.Q),
			param.Offset+param.Limit, param.Limit)
	}

	list := QuerySearchList{
		Total:  int(total),
		Offset: param.Offset,
		Limit:  param.Limit,
		Prev:   prev,
		Next:   next,
		List:   chars,
	}

	return c.JSON(http.StatusOK, list)
}

// GetCharSearch - get char search
// @Summary search a character
// @Description get a character list from a character
// @Produce json
// @Param q query string true "a single character"
// @Param offset query int false "default: 0"
// @Param limit query int false "default: 10; '0' means default; '-1' means no limitation"
// @Param delegate query bool false "if 'true', pick up a single character for each title"
// @Success 200 {object} mojiportal.List
// @Failure 400 {object} mojiportal.ErrorResponse
// @Failure 500 {object} mojiportal.ErrorResponse
// @Router /char/search [get]
func GetCharSearch(c echo.Context) error {
	errRes := new(mojiportal.ErrorResponse)

	param := new(mojiportal.CharSearchParam)
	if err := c.Bind(param); err != nil {
		errRes.Error = append(errRes.Error, mojiportal.ErrorItem{
			Key:     "query",
			Message: err.(*echo.HTTPError).Message.(string),
		})
		return c.JSON(http.StatusBadRequest, errRes)
	}

	if err := c.Validate(param); err != nil {
		return c.JSON(http.StatusBadRequest,
			mojiportal.GetErrorResponseFromValidateError(err))
	}

	param.Unicode = int([]rune(param.Q)[0])

	if param.Limit == 0 {
		param.Limit = 10
	}

	var jikeis []JikeiPageView

	sql := "SELECT " + jpvKeys + " FROM jikei AS j INNER JOIN page AS p ON j.pid = p.pid WHERE unicode = :Unicode ORDER BY jid"
	if param.Offset != 0 && !param.OfEach {
		sql += " OFFSET :Offset"
	}
	if param.Limit != -1 && !param.OfEach {
		sql += " LIMIT :Limit"
	}
	_, err := dbmap.Select(&jikeis, sql, param)
	if err != nil {
		errRes.Error = append(errRes.Error, mojiportal.ErrorItem{
			Key:     "query",
			Message: err.Error(),
		})
		return c.JSON(http.StatusBadRequest, errRes)
	}

	chars := make([]mojiportal.Char, 0, len(jikeis))
	delegate := true
	var prevBID string

	for _, j := range jikeis {
		if j.BID == prevBID {
			if param.OfEach {
				continue
			} else {
				delegate = false
			}
		} else {
			delegate = true
			prevBID = j.BID
		}
		ch := string(j.Unicode)
		chars = append(chars, mojiportal.Char{
			ID:        j.JID,
			Character: ch,
			Delegate:  delegate,
			Unicode:   fmt.Sprintf("%04X", j.Unicode),
			Source: mojiportal.Source{
				Label: fmt.Sprintf("『%s』（%s）%d%s",
					j.Title, j.BID, j.Frame, j.Side),
				BID:   j.BID,
				Title: j.Title,
				Frame: j.Frame,
				Side:  j.Side,
			},
			ThumbnailURL: j.GetIIIFImageID() + "/full/full/0/default.jpg",
			ManifestURL:  baseURL + "/api/char/" + j.JID + "/manifest",
			Link: fmt.Sprintf("%s/viewer/?manifest=%s/api/biblio/%s/%s/manifest&pos=%d&xywh=%d,%d,%d,%d",
				baseURL, baseURL, j.BID, url.PathEscape(ch),
				j.Pos, j.X, j.Y, j.Width, j.Height),
			Subject:   "日本古典籍くずし字データ",
			Creator:   "国文学研究資料館",
			Rights:    "『日本古典籍くずし字データセット』（国文研ほか所蔵／CODH加工） doi:10.20676/00000340",
			RightsURL: "https://creativecommons.org/licenses/by-sa/4.0/deed.ja",
		})
	}

	if param.OfEach {
		if param.Offset >= len(chars) {
			chars = chars[0:0]
		} else if param.Limit == -1 {
			chars = chars[param.Offset:]
		} else if param.Offset+param.Limit > len(chars) {
			chars = chars[param.Offset:]
		} else {
			chars = chars[param.Offset : param.Offset+param.Limit]
		}
	}

	list := mojiportal.List{
		SearchResults: len(chars),
		List:          chars,
	}

	return c.JSON(http.StatusOK, list)
}

// PostPageListUpload - post page list upload
func PostPageListUpload(c echo.Context) error {
	titles, err := GetTitlesFromCODH()
	if err != nil {
		return err
	}

	// read csv
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	reader := csv.NewReader(src)
	var cols []string
	created := time.Now().Format("2006-01-02 15:04:05")
	cntTotal := 0
	cntInserted := 0
	cntUpdated := 0
	cntError := 0

	prevBID := ""
	pos := 1

	// ignore the first line
	_, _ = reader.Read()

	for {
		cols, err = reader.Read()
		if err != nil {
			break
		}

		cntTotal++

		frame, err := strconv.Atoi(cols[2])
		if err != nil {
			log.Printf("ERROR: %v", err)
			cntError++
			continue
		}

		side := ""
		if cols[3] == "1" {
			side = "r"
		} else if cols[3] == "2" {
			side = "v"
		}

		width, err := strconv.Atoi(cols[4])
		if err != nil {
			log.Printf("ERROR: %v", err)
			cntError++
			continue
		}

		height, err := strconv.Atoi(cols[5])
		if err != nil {
			log.Printf("ERROR: %v", err)
			cntError++
			continue
		}

		if cols[1] == prevBID {
			pos++
		} else {
			prevBID = cols[1]
			pos = 1
		}

		title, _ := titles[cols[1]]

		page := &Page{
			PID:     cols[0],
			BID:     cols[1],
			Title:   title,
			Pos:     pos,
			Frame:   frame,
			Side:    side,
			Width:   width,
			Height:  height,
			Created: created,
		}
		log.Printf("INFO: %d: %v", cntTotal, page)

		err = dbmap.Insert(page)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") ||
				strings.Contains(err.Error(), "一意性制約") {
				i, err := dbmap.Update(page)
				if err != nil || i != 1 {
					log.Printf("Error: update error: %d, %s", i, err)
					cntError++
					continue
				} else {
					cntUpdated++
					log.Printf("INFO: successfully updated")
				}
			} else {
				log.Printf("ERROR: %s", err)
				cntError++
				continue
			}
		} else {
			cntInserted++
			log.Printf("INFO: successfully inserted")
		}
	}

	log.Printf("INFO: total: %d; inserted %d; updated: %d, error: %d",
		cntTotal, cntInserted, cntUpdated, cntError)
	return c.JSON(http.StatusOK, ResponseListUpload{
		Message:  "successfully registered",
		Total:    cntTotal,
		Inserted: cntInserted,
		Updated:  cntUpdated,
		Error:    cntError,
	})
}

// PostJikeiListUpload - post jikei list upload
func PostJikeiListUpload(c echo.Context) error {
	// read csv
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	reader := csv.NewReader(src)
	var cols []string
	created := time.Now().Format("2006-01-02 15:04:05")
	cntTotal := 0
	cntInserted := 0
	cntUpdated := 0
	cntError := 0

	pages := map[string]Page{}

	// ignore the first line
	_, _ = reader.Read()

	for {
		cols, err = reader.Read()
		if err != nil {
			break
		}

		cntTotal++

		unicode, err := strconv.ParseInt(cols[0][2:], 16, 0)
		if err != nil {
			log.Printf("ERROR: %v", err)
			cntError++
			continue
		}

		pid := cols[1]

		x, err := strconv.Atoi(cols[2])
		if err != nil {
			log.Printf("ERROR: %v", err)
			cntError++
			continue
		}

		y, err := strconv.Atoi(cols[3])
		if err != nil {
			log.Printf("ERROR: %v", err)
			cntError++
			continue
		}

		width, err := strconv.Atoi(cols[6])
		if err != nil {
			log.Printf("ERROR: %v", err)
			cntError++
			continue
		}

		height, err := strconv.Atoi(cols[7])
		if err != nil {
			log.Printf("ERROR: %v", err)
			cntError++
			continue
		}

		page, ok := pages[pid]
		if !ok {
			err = dbmap.SelectOne(&page,
				"SELECT * FROM page WHERE pid = $1", pid)
			if err != nil {
				log.Printf("%v; page %s is not registered yet\n", err, pid)
				cntError++
				continue
			}
			pages[pid] = page
		}

		jid := fmt.Sprintf("%s_%s_X%04d_Y%04d", cols[0], pid, x, y)

		// title, _ := titles[bid]

		j := &Jikei{
			JID:     jid,
			Unicode: int(unicode),
			PID:     pid,
			X:       x,
			Y:       y,
			BlockID: cols[4],
			CharID:  cols[5],
			Width:   width,
			Height:  height,
			Created: created,
		}
		log.Printf("INFO: %d: %v", cntTotal, j)

		err = dbmap.Insert(j)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") ||
				strings.Contains(err.Error(), "一意性制約") {
				i, err := dbmap.Update(j)
				if err != nil || i != 1 {
					log.Printf("Error: update error: %d, %s", i, err)
					cntError++
					continue
				} else {
					cntUpdated++
					log.Printf("INFO: successfully updated")
				}
			} else {
				log.Printf("ERROR: %s", err)
				cntError++
				continue
			}
		} else {
			cntInserted++
			log.Printf("INFO: successfully inserted")
		}
	}

	log.Printf("INFO: total: %d; inserted %d; updated: %d, error: %d",
		cntTotal, cntInserted, cntUpdated, cntError)
	return c.JSON(http.StatusOK, ResponseListUpload{
		Message:  "successfully registered",
		Total:    cntTotal,
		Inserted: cntInserted,
		Updated:  cntUpdated,
		Error:    cntError,
	})
}

/* utilities */

// FormatManifest - format manifest
func FormatManifest(pages []Page, param *QueryManifestParam) iiif.Manifest {
	chars := param.Q

	canvases := make([]iiif.Canvas, 0, len(pages))
	for _, page := range pages {
		var canvas string
		var otherContents []iiif.OtherContent
		if chars == "" {
			canvas = fmt.Sprintf("%s/api/page/%s/canvas/c1", baseURL, page.PID)
		} else {
			canvas = fmt.Sprintf("%s/api/page/%s/%s/canvas/c1",
				baseURL, page.PID, url.PathEscape(chars))
			otherContents = append(otherContents, iiif.OtherContent{
				ID:   canvas + "/annolist",
				Type: "sc:AnnotationList",
			})
		}

		jpg := page.GetIIIFImageID()

		canvases = append(canvases, iiif.Canvas{
			ID:   canvas,
			Type: "sc:Canvas",
			Label: fmt.Sprintf("%s %d%s",
				page.Title, page.Frame, page.Side),
			Width:  page.Width,
			Height: page.Height,
			Images: []iiif.Image{
				iiif.Image{
					ID:         canvas + "/annotion/anno1",
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
					On: canvas,
				},
			},
			OtherContent: otherContents,
		})
	}

	manifest := iiif.Manifest{
		Context: "http://iiif.io/api/presentation/2/context.json",
		ID: fmt.Sprintf("%s/api/manifest?q=%s",
			baseURL, url.QueryEscape(chars)),
		Type:             "sc:Manifest",
		Label:            chars,
		ViewingHint:      "paged",
		ViewingDirection: "right-to-left",
		License:          manifestLicense,
		Attribution:      manifestAttribution,
		Logo:             baseURL + "/img/nijl_symbolmark.jpg",
		Related: []iiif.IDFormat{
			iiif.IDFormat{
				ID:     "https://kotenseki.nijl.ac.jp/",
				Format: "text/html",
			},
		},
		Within: baseURL + "/",
		Sequences: []iiif.Sequence{
			iiif.Sequence{
				ID: fmt.Sprintf("%s/api/sequence?q=%s",
					baseURL, url.QueryEscape(chars)),
				Type:     "sc:Sequence",
				Canvases: canvases,
			},
		},
	}

	return manifest
}

// FormatBiblioManifest - format biblio manifest
func FormatBiblioManifest(pages []Page, chars string) iiif.Manifest {
	canvases := make([]iiif.Canvas, 0, len(pages))
	for _, page := range pages {
		var canvas string
		var otherContents []iiif.OtherContent
		if chars == "" {
			canvas = fmt.Sprintf("%s/api/page/%s/canvas/c1", baseURL, page.PID)
		} else {
			canvas = fmt.Sprintf("%s/api/page/%s/%s/canvas/c1",
				baseURL, page.PID, url.PathEscape(chars))
			otherContents = append(otherContents, iiif.OtherContent{
				ID:   canvas + "/annolist",
				Type: "sc:AnnotationList",
			})
		}

		jpg := page.GetIIIFImageID()

		canvases = append(canvases, iiif.Canvas{
			ID:     canvas,
			Type:   "sc:Canvas",
			Label:  fmt.Sprintf("%d%s", page.Frame, page.Side),
			Width:  page.Width,
			Height: page.Height,
			Images: []iiif.Image{
				iiif.Image{
					ID:         canvas + "/annotion/anno1",
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
					On: canvas,
				},
			},
			OtherContent: otherContents,
		})
	}

	base := baseURL + "/api/biblio/" + pages[0].BID
	if chars != "" {
		base += "/" + url.PathEscape(chars)
	}

	manifest := iiif.Manifest{
		Context: "http://iiif.io/api/presentation/2/context.json",
		ID:      base + "/manifest",
		Type:    "sc:Manifest",
		Label:   fmt.Sprintf("『%s』（%s）", pages[0].Title, pages[0].BID),
		Metadata: []iiif.Metadatum{
			iiif.Metadatum{
				Label: "BID",
				Value: pages[0].BID,
			},
			iiif.Metadatum{
				Label: "TITLE",
				Value: pages[0].Title,
			},
		},
		ViewingHint:      "paged",
		ViewingDirection: "right-to-left",
		License:          manifestLicense,
		Attribution:      manifestAttribution,
		Logo:             baseURL + "/img/nijl_symbolmark.jpg",
		Related: []iiif.IDFormat{
			iiif.IDFormat{
				ID:     fmt.Sprintf("https://kotenseki.nijl.ac.jp/biblio/%s/viewer", pages[0].BID),
				Format: "text/html",
			},
		},
		Within: baseURL + "/",
		Sequences: []iiif.Sequence{
			iiif.Sequence{
				ID:       base + "/sequence",
				Type:     "sc:Sequence",
				Canvases: canvases,
			},
		},
	}

	return manifest
}

// FormatPageAnnoList - format page annolist
func FormatPageAnnoList(jikeis []JikeiPageView, image, chars string) iiif.AnnoList {
	var canvas string
	if chars == "" {
		canvas = fmt.Sprintf("%s/api/page/%s/canvas/c1", baseURL, image)
	} else {
		canvas = fmt.Sprintf("%s/api/page/%s/%s/canvas/c1",
			baseURL, image, url.PathEscape(chars))
	}
	annoResources := make([]iiif.AnnoResources, 0, len(jikeis))
	for _, j := range jikeis {
		annoResources = append(annoResources, iiif.AnnoResources{
			ID:         fmt.Sprintf("%s/annolist/%s", canvas, j.JID),
			Context:    "http://iiif.io/api/presentation/2/context.json",
			Type:       "oa:Annotation",
			Motivation: []string{"oa:commenting"},
			Resource: []iiif.AnnoResource{
				iiif.AnnoResource{
					Type:   "dctypes:Text",
					Format: "text/html",
					Chars: fmt.Sprintf("%s xywh=%d,%d,%d,%d",
						string(j.Unicode),
						j.X, j.Y, j.Width, j.Height),
				},
			},
			On: fmt.Sprintf("%s#xywh=%d,%d,%d,%d",
				canvas, j.X, j.Y, j.Width, j.Height),
		})
	}

	annoList := iiif.AnnoList{
		ID:        fmt.Sprintf("%s/annolist", canvas),
		Context:   "http://www.shared-canvas.org/ns/context.json",
		Type:      "sc:AnnotationList",
		Resources: annoResources,
	}
	return annoList
}

// FormatCharManifest - format character manifest
func FormatCharManifest(j JikeiPageView) iiif.Manifest {
	base := baseURL + "/api/char/" + j.JID
	jpg := j.GetIIIFImageID()
	ch := string(j.Unicode)
	manifest := iiif.Manifest{
		Context: "http://iiif.io/api/presentation/2/context.json",
		ID:      base + "/manifest",
		Type:    "sc:Manifest",
		Label: fmt.Sprintf("「%s」『%s』（%s）%d%s",
			ch, j.Title, j.BID, j.Frame, j.Side),
		Metadata: []iiif.Metadatum{
			iiif.Metadatum{
				Label: "CHARACTER",
				Value: ch,
			},
			iiif.Metadatum{
				Label: "CHARACTER CODE",
				Value: fmt.Sprintf("%X", j.Unicode),
			},
			iiif.Metadatum{
				Label: "BID",
				Value: j.BID,
			},
			iiif.Metadatum{
				Label: "TITLE",
				Value: j.Title,
			},
			iiif.Metadatum{
				Label: "FRAME",
				Value: strconv.Itoa(j.Frame),
			},
			iiif.Metadatum{
				Label: "SIDE",
				Value: j.Side,
			},
			iiif.Metadatum{
				Label: "X",
				Value: strconv.Itoa(j.X),
			},
			iiif.Metadatum{
				Label: "Y",
				Value: strconv.Itoa(j.Y),
			},
			iiif.Metadatum{
				Label: "WIDTH",
				Value: strconv.Itoa(j.Width),
			},
			iiif.Metadatum{
				Label: "HEIGHT",
				Value: strconv.Itoa(j.Height),
			},
		},
		ViewingHint: "paged",
		License:     manifestLicense,
		Attribution: manifestAttribution,
		Logo:        baseURL + "/img/nijl_symbolmark.jpg",
		Related: []iiif.IDFormat{
			iiif.IDFormat{
				ID:     "http://codh.rois.ac.jp/char-shape/",
				Format: "text/html",
			},
			iiif.IDFormat{
				ID: fmt.Sprintf("%s/viewer/?manifest=%s/api/biblio/%s/%s/manifest&pos=%d&xywh=%d,%d,%d,%d",
					baseURL, baseURL, j.BID, url.PathEscape(ch),
					j.Pos, j.X, j.Y, j.Width, j.Height),
				Format: "text/html",
			},
		},
		Within: baseURL + "/",
		Sequences: []iiif.Sequence{
			iiif.Sequence{
				ID:   base + "/sequence",
				Type: "sc:Sequence",
				Canvases: []iiif.Canvas{
					iiif.Canvas{
						ID:     base + "/canvas/c1",
						Type:   "sc:Canvas",
						Label:  ch,
						Width:  j.Width,
						Height: j.Height,
						Images: []iiif.Image{
							iiif.Image{
								ID:         base + "/canvas/c1/annotion/anno1",
								Type:       "oa:Annotation",
								Motivation: "sc:painting",
								Resource: iiif.Resource{
									ID:     jpg + "/full/full/0/default.jpg",
									Type:   "dctypes:Image",
									Format: "image/jpeg",
									Width:  j.Width,
									Height: j.Height,
									Service: iiif.Service{
										Context: "http://iiif.io/api/image/2/context.json",
										ID:      jpg,
										Profile: "http://iiif.io/api/image/2/level1.json",
									},
								},
								On: base + "/canvas/c1",
							},
						},
					},
				},
			},
		},
	}
	manifest.Context = "http://iiif.io/api/presentation/2/context.json"
	return manifest
}
