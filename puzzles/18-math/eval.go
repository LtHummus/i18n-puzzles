package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type RPNTokenKind int

const (
	TokenKindOperator RPNTokenKind = iota
	TokenKindOperand
)

type RPNToken struct {
	Value any
	Kind  RPNTokenKind
}

func (rt *RPNToken) isOperand() bool {
	return rt.Kind == TokenKindOperand
}

func (rt *RPNToken) numericValue() float64 {
	if rt.Kind != TokenKindOperand {
		panic("attempted to get number from operator")
	}

	return rt.Value.(float64)
}

func (rt *RPNToken) operatorValue() string {
	if rt.Kind != TokenKindOperator {
		panic("attempted to get operator from operand")
	}

	return rt.Value.(string)
}

func (rt *RPNToken) String() string {
	return fmt.Sprintf("%v", rt.Value)
}

func newNumberToken(x float64) *RPNToken {
	return &RPNToken{
		Kind:  TokenKindOperand,
		Value: x,
	}
}

func newOperatorToken(x string) *RPNToken {
	if x != "+" && x != "-" && x != "*" && x != "/" && x != "(" && x != ")" {
		panic("invalid operator")
	}
	
	return &RPNToken{
		Kind:  TokenKindOperator,
		Value: x,
	}
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "^":
		return 3
	}
	return 0
}

func shuntingYard(tokens []string) []*RPNToken {
	var output []*RPNToken
	var operators []string

	for _, token := range tokens {
		if value, err := strconv.ParseFloat(token, 64); err == nil {
			// token is a number, so add it to the output
			output = append(output, newNumberToken(value))
		} else if token == "(" {
			// open paren has been found, so add that too
			operators = append(operators, token)
		} else if token == ")" {
			// find the matching close
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, newOperatorToken(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			if len(operators) > 0 {
				operators = operators[:len(operators)-1]
			}
		} else {
			// operator found, figure out it's precednce and act accordingly
			for len(operators) > 0 &&
				operators[len(operators)-1] != "(" &&
				(precedence(operators[len(operators)-1]) > precedence(token) ||
					(precedence(operators[len(operators)-1]) == precedence(token))) {
				output = append(output, newOperatorToken(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		}
	}

	for len(operators) > 0 {
		output = append(output, newOperatorToken(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	return output
}

func tokenize(input string) ([]string, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))

	var tok rune
	var result = make([]string, 0)
	for tok != scanner.EOF {
		tok = s.Scan()
		value := strings.TrimSpace(s.TokenText())
		if len(value) > 0 {
			result = append(result, s.TokenText())
		}
	}
	return result, nil
}

func evalOperator(op string, a float64, b float64) float64 {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	default:
		panic("unknown operator")
	}
}

func evalExpression(x string) float64 {
	tokens, err := tokenize(x)
	if err != nil {
		panic(err)
	}

	rpn := shuntingYard(tokens)

	var stack []float64
	for _, token := range rpn {
		if token.isOperand() {
			stack = append(stack, token.numericValue())
		} else {
			if len(stack) < 2 {
				panic("need 2 operands")
			}
			a, b := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			evalRes := evalOperator(token.operatorValue(), a, b)
			stack = append(stack, evalRes)
		}
	}

	if len(stack) != 1 {
		panic("invalid number of items on stack at the end")
	}

	return stack[0]
}
