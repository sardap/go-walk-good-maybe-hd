package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	path := os.Args[1]

	r := gin.Default()

	r.Use(cors.Default())

	r.StaticFS("/assets", http.Dir(path))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
