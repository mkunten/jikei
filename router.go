package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// _ "github.com/mkunten/jikei/docs"
	// "github.com/swaggo/echo-swagger"
)

// Router - echo routing
func Router(e *echo.Echo) {
	// // swaggo
	// e.GET("/jikei/api/docs/*", echoSwagger.WrapHandler)

	// groups
	api := e.Group("/jikei/api")
	api.Use(middleware.Gzip())
	admin := api.Group("/admin")

	// for ordinal user
	api.GET("/jikeis", func(c echo.Context) error {
		var jikeis []Jikei
		_, err := dbmap.Select(&jikeis, "SELECT * FROM jikei ORDER BY id desc LIMIT 10")
		if err != nil {
			c.Logger().Error(" Select: ", err)
			return c.String(http.StatusBadRequest, "Select: "+err.Error())
		}
		return c.JSON(http.StatusOK, jikeis)
	})
	api.GET("/biblio/:bid/manifest", GetBiblioManifest)
	api.GET("/biblio/:bid/:chars/manifest", GetBiblioManifest)
	api.GET("/page/:pid/:chars/canvas/c1/annolist", GetPageAnnoList)
	api.GET("/page/:pid/manifest", GetPageManifest)
	api.GET("/page/:pid/:chars/manifest", GetPageManifest)
	api.GET("/page/:pid/canvas/c1/annolist", GetPageAnnoList)
	api.GET("/char/:jid/manifest", GetCharManifest)
	api.GET("/char/search", GetCharSearch)
	api.GET("/manifest", GetQueryManifest)
	api.GET("/search", GetQuerySearch)

	// for administrator
	admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "admin" && password == "jikei-admin" {
			return true, nil
		}
		return false, nil
	}))
	admin.POST("/jikeilistupload", PostJikeiListUpload)
	admin.POST("/pagelistupload", PostPageListUpload)
}
