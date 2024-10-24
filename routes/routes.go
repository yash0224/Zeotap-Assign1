package routes

import (
	"zeotap_assign1/controllers"

	"github.com/gin-gonic/gin"
)

var RegisterRoutes = func(c *gin.Engine) {
	c.POST("/create", controllers.CreateRule)
	c.POST("/combine", controllers.CombineRules)
	c.POST("/evaluate", controllers.EvaluateRule)
}
