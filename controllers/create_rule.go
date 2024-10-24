package controllers

import (
	"context"
	"fmt"
	"net/http"

	"zeotap_assign1/database"

	"github.com/gin-gonic/gin"
)

func CreateRule(c *gin.Context) {
	var body struct {
		RuleID     string `json:"rule_id"`
		RuleString string `json:"rule_string"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fmt.Printf("The rule string is: %s", body.RuleString)

	ast, err := parseRuleString(body.RuleString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule string"})
		return
	}

	newRule := &database.Rule{
		RuleString: body.RuleString,
		AST:        ast,
		RuleID:     body.RuleID,
	}
	collection := database.GetCollection(database.Client, "rule")
	collection.InsertOne(context.TODO(), newRule)

	c.JSON(http.StatusCreated, newRule)
}
