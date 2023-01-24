package web

import (
	"NUMParser/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func setupRouter() *gin.Engine {
	//gin.DisableConsoleColor()
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

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

func Start(port string) {
	go func() {
		r := setupRouter()
		err := r.Run(":" + port)
		if err != nil {
			log.Println("Error start web server on port", port, ":", err)
		}
	}()
}
