package controllers

import (
	"errors"
	"fmt"
	"strings"

	"zeotap_assign1/model"
)

func parseRuleString(ruleString string) (*model.Node, error) {
	fmt.Println("Parsing rule string:", ruleString)

	// Pre-process the rule string to handle parentheses properly
	ruleString = strings.ReplaceAll(ruleString, "(", " ( ")
	ruleString = strings.ReplaceAll(ruleString, ")", " ) ")
	tokens := strings.Fields(ruleString)

	if len(tokens) == 0 {
		return nil, errors.New("invalid rule string")
	}

	outputQueue := []string{}
	operatorStack := []string{}
	precedence := map[string]int{
		"AND": 2,
		"OR":  1,
		"(":   0,
	}

	for _, token := range tokens {
		switch token {
		case "AND", "OR":
			for len(operatorStack) > 0 &&
				operatorStack[len(operatorStack)-1] != "(" &&
				precedence[operatorStack[len(operatorStack)-1]] >= precedence[token] {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)
		case "(":
			operatorStack = append(operatorStack, token)
		case ")":
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] == "(" {
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
		default:
			outputQueue = append(outputQueue, token)
		}
	}

	for len(operatorStack) > 0 {
		if operatorStack[len(operatorStack)-1] == "(" {
			return nil, errors.New("mismatched parentheses")
		}
		outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	fmt.Println("Output queue =", outputQueue)

	var stack []*model.Node

	createOperandNode := func(attribute, operator, value string) *model.Node {
		return &model.Node{
			Type: "operand",
			Value: map[string]interface{}{
				"attribute": attribute,
				"operator":  operator,
				"value":     value,
			},
		}
	}

	for i := 0; i < len(outputQueue); i++ {
		token := outputQueue[i]
		switch token {
		case "AND", "OR":
			if len(stack) < 2 {
				return nil, errors.New("invalid expression: not enough operands")
			}
			right := stack[len(stack)-1]
			left := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, &model.Node{
				Type:  "operator",
				Left:  left,
				Right: right,
				Value: token,
			})
		case ">", "<", "=", ">=", "<=":
			if len(stack) < 1 || i >= len(outputQueue)-1 {
				return nil, errors.New("invalid expression: incomplete comparison")
			}
			value := outputQueue[i+1]
			attribute := stack[len(stack)-1].Value.(string)
			stack = stack[:len(stack)-1]
			stack = append(stack, createOperandNode(attribute, token, value))
			i++
		default:
			stack = append(stack, &model.Node{
				Type:  "operand",
				Value: token,
			})
		}
	}

	if len(stack) != 1 {
		return nil, errors.New("invalid expression: incorrect number of operands")
	}

	return stack[0], nil
}
