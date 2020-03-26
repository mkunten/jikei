package main

// @title 国文研字形検索βAPI
// @version 0.0.1

// @license.name MIT

// @host lab.nijl.ac.jp
// @BasePath /jikei/api
// @schemes https

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/gorp.v2"
)

const (
	dsn         = "user=jikei password=jikei-admin dbname=jikei sslmode=disable"
	iiifBaseURL = "https://lab.nijl.ac.jp/jikei/iiif"
	imgRoot     = "/opt/loris2/images"

	manifestLicense     = "https://creativecommons.org/licenses/by-sa/4.0/deed.ja"
	manifestAttribution = "国文学研究資料館/『日本古典籍くずし字データセット』（国文研ほか所蔵／CODH加工） doi:10.20676/00000340"
)

var (
	version string // -ldflags "-X main.version=<version>"
	baseURL = "https://lab.nijl.ac.jp/jikei"
	dbmap   *gorp.DbMap
)

func main() {
	if version == "" {
		baseURL = "http://localhost:58080/jikei"
	} else {
		println("version: " + version)
	}

	/* initialize db */
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbmap = &gorp.DbMap{
		Db:      db,
		Dialect: gorp.PostgresDialect{},
	}

	t := dbmap.AddTableWithName(Jikei{}, "jikei")
	t.AddIndex("jikei_pid_idx", "Btree", []string{"pid"})
	t.AddIndex("jikei_unicode_idx", "Btree", []string{"unicode"})

	t = dbmap.AddTableWithName(Page{}, "page")
	t.AddIndex("page_bid_idx", "Btree", []string{"bid"})

	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "jikei:", log.Lmicroseconds))

	if err = dbmap.CreateTablesIfNotExists(); err != nil {
		log.Fatal(err)
	}

	if err = dbmap.CreateIndex(); err != nil &&
		!strings.HasSuffix(err.Error(), "already exists") &&
		!strings.HasSuffix(err.Error(), "すでに存在します") {
		log.Fatal(err)
	}

	/* initialize echo */
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Validator = &CustomValidator{validator: validator.New()}

	/* routing */
	Router(e)

	/* echo start */
	e.Logger.Fatal(e.Start(":8080"))
}
