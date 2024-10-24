package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"zeotap_assign1/database"
	"zeotap_assign1/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func evaluateNode(node *model.Node, data map[string]interface{}) (bool, error) {
	if node.Type == "operator" {
		leftValue, err := evaluateNode(node.Left, data)
		if err != nil {
			return false, err
		}
		rightValue, err := evaluateNode(node.Right, data)
		if err != nil {
			return false, err
		}

		switch node.Value {
		case "AND":
			return leftValue && rightValue, nil
		case "OR":
			return leftValue || rightValue, nil
		default:
			return false, errors.New("invalid operator")
		}
	} else if node.Type == "operand" {
		valueMap := node.Value.(map[string]interface{})
		attribute := valueMap["attribute"].(string)
		operator := valueMap["operator"].(string)
		cleanValue := valueMap["value"]

		dataValue, ok := data[attribute]
		if !ok {
			return false, errors.New("attribute not found in data")
		}

		switch operator {
		case ">":
			return dataValue.(float64) > cleanValue.(float64), nil
		case "<":
			return dataValue.(float64) < cleanValue.(float64), nil
		case ">=":
			return dataValue.(float64) >= cleanValue.(float64), nil
		case "<=":
			return dataValue.(float64) <= cleanValue.(float64), nil
		case "=":
			return dataValue == cleanValue, nil
		case "!=":
			return dataValue != cleanValue, nil
		default:
			return false, errors.New("invalid comparison operator")
		}
	}

	return false, errors.New("invalid node type")
}

func EvaluateRule(c *gin.Context) {
	var body struct {
		RuleID string                 `json:"ruleId"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var ast *model.Node
	collection := database.GetCollection(database.Client, "rule")
	err := collection.FindOne(context.TODO(), bson.M{"ruleid": body.RuleID}).Decode(&ast)

	fmt.Println(ast)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"message": "rule not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := evaluateNode(ast, body.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
