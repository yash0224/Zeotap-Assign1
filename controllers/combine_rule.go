package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"zeotap_assign1/database"
	"zeotap_assign1/model"

	"github.com/gin-gonic/gin"
)

// operatorPrecedence defines the precedence of logical operators
var operatorPrecedence = map[string]int{
	"AND": 2,
	"OR":  1,
}

func CombineRules(c *gin.Context) {
	var body struct {
		RuleID          string   `json:"rule_id"`
		RuleStrings     []string `json:"ruleStrings"`
		CombineOperator string   `json:"combop"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var asts []*model.Node
	for _, ruleString := range body.RuleStrings {
		ast, err := parseRuleString(ruleString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing rule: " + err.Error()})
			return
		}
		asts = append(asts, ast)
	}

	combinedAST := combineASTsHeuristically(asts, strings.ToUpper(body.CombineOperator))
	optimizedAST := optimizeAST(combinedAST)

	newRule := &database.Rule{
		RuleID:     body.RuleID,
		RuleString: generateRuleString(optimizedAST),
		AST:        optimizedAST,
	}

	collection := database.GetCollection(database.Client, "rule")
	_, err := collection.InsertOne(context.TODO(), newRule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"combinedAST": optimizedAST,
		"ruleString":  newRule.RuleString,
	})
}

func combineASTsHeuristically(asts []*model.Node, operator string) *model.Node {
	if len(asts) == 0 {
		return nil
	}
	if len(asts) == 1 {
		return asts[0]
	}

	result := asts[0]
	for i := 1; i < len(asts); i++ {
		result = &model.Node{
			Type:  "operator",
			Value: operator,
			Left:  result,
			Right: asts[i],
		}
	}
	return result
}

// func optimizeAST(node *model.Node) *model.Node {
// 	if node == nil {
// 		return nil
// 	}

// 	// Base case: leaf nodes
// 	if node.Type != "operator" {
// 		return node
// 	}

// 	// Recursively optimize children
// 	node.Left = optimizeAST(node.Left)
// 	node.Right = optimizeAST(node.Right)

// 	return node
// }

func generateRuleString(node *model.Node) string {
	if node == nil {
		return ""
	}

	switch node.Type {
	case "operator":
		left := generateRuleString(node.Left)
		right := generateRuleString(node.Right)
		operator, ok := node.Value.(string)
		if !ok {
			return fmt.Sprintf("(%s INVALID_OPERATOR %s)", left, right)
		}
		return fmt.Sprintf("(%s %s %s)", left, operator, right)

	case "operand":
		switch v := node.Value.(type) {
		case string:
			return v
		case map[string]interface{}:
			attribute, _ := v["attribute"].(string)
			operator, _ := v["operator"].(string)
			value, _ := v["value"].(string)
			// Remove quotes from string values except for string literals
			if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
				return fmt.Sprintf("%s %s %s", attribute, operator, value)
			}
			return fmt.Sprintf("%s %s %s", attribute, operator, value)
		default:
			return "INVALID_VALUE"
		}
	default:
		return "INVALID_NODE_TYPE"
	}
}
func optimizeAST(node *model.Node) *model.Node {
	if node == nil {
		return nil
	}

	// Base case: leaf nodes
	if node.Type != "operator" {
		return node
	}

	// Recursively optimize children
	node.Left = optimizeAST(node.Left)
	node.Right = optimizeAST(node.Right)

	// If it's an AND operation with identical subtrees, return just one subtree
	if node.Type == "operator" {
		if operator, ok := node.Value.(string); ok && operator == "AND" {
			if areNodesEqual(node.Left, node.Right) {
				return node.Left
			}
		}
	}

	return node
}

// areNodesEqual compares two AST nodes for structural equality
func areNodesEqual(node1, node2 *model.Node) bool {
	if node1 == nil && node2 == nil {
		return true
	}
	if node1 == nil || node2 == nil {
		return false
	}
	if node1.Type != node2.Type {
		return false
	}

	// Compare values based on node type
	switch node1.Type {
	case "operator":
		if node1.Value != node2.Value {
			return false
		}
		return areNodesEqual(node1.Left, node2.Left) && areNodesEqual(node1.Right, node2.Right)
	case "operand":
		// For operand nodes, compare the actual values
		switch v1 := node1.Value.(type) {
		case string:
			v2, ok := node2.Value.(string)
			return ok && v1 == v2
		case map[string]interface{}:
			v2, ok := node2.Value.(map[string]interface{})
			if !ok {
				return false
			}
			return v1["attribute"] == v2["attribute"] &&
				v1["operator"] == v2["operator"] &&
				v1["value"] == v2["value"]
		default:
			return false
		}
	default:
		return false
	}
}
