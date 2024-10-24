package model

type Node struct {
	Type  string      `json:"type"` // type operator or operand
	Left  *Node       `json:"left"` // letf child
	Right *Node       `json:"right"` // right child
	Value interface{} `json:"value"` // value of this code
}
