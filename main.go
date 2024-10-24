package main

import (
	"net/http"

	"zeotap_assign1/database"
	"zeotap_assign1/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func init() {
	godotenv.Load()
	database.InitializeConnections()
}

func main() {
	r := gin.Default()

	r.LoadHTMLFiles("./templates/index.html")

	routes.RegisterRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Wrap your router with the CORS middleware
	handler := c.Handler(r)

	// Start the server
	http.ListenAndServe(":8080", handler)
}
