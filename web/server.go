package web

import (
	"NUMParser/config"
	"NUMParser/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var route *gin.Engine

func setupRouter() *gin.Engine {
	//gin.DisableConsoleColor()
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.Static("/css", "public/css")
	r.Static("/img", "public/img")
	r.Static("/js", "public/js")
	r.StaticFile("/", "public/index.html")

	// http://127.0.0.1:38888/search?query=venom
	r.GET("/search", func(c *gin.Context) {
		if query, ok := c.GetQuery("query"); ok {
			torrs := db.SearchTorr(query)
			c.JSON(200, torrs)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	})

	return r
}

var isSetStatic bool

func SetStaticReleases() {
	if !isSetStatic {
		route.Static("/releases", config.SaveReleasePath)
		isSetStatic = true
	}
}

func Start(port string) {
	go func() {
		route = setupRouter()
		err := route.Run(":" + port)
		if err != nil {
			log.Println("Error start web server on port", port, ":", err)
		}
	}()
}
