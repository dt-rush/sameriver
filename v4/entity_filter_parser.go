package sameriver

import (
	"errors"
	"fmt"
	"strings"
)

/*
grammar:

Expr            := PredicateExpr (Semicolon SortExpr)?
PredicateExpr   := Not? Function (And PredicateExpr | Or PredicateExpr)?
Function        := Identifier OpenParen Args CloseParen
Args            := Identifier (Comma Identifier)*
*/

type NodeType int

const (
	NodeExpr NodeType = iota
	NodePredicateExpr
	NodeSortExpr
	NodeNot
	NodeAnd
	NodeOr
	NodeFunction
	NodeIdentifier
)

var nodeTypeStrings = map[NodeType]string{
	NodeExpr:          "NodeExpr",
	NodePredicateExpr: "NodePredicateExpr",
	NodeSortExpr:      "NodeSortExpr",
	NodeNot:           "NodeNot",
	NodeAnd:           "NodeAnd",
	NodeOr:            "NodeOr",
	NodeFunction:      "NodeFunction",
	NodeIdentifier:    "NodeIdentifier",
}

type Node struct {
	Type     NodeType
	Value    string
	Children []*Node
}

func (n *Node) String() string {
	chStr := ""
	for i, ch := range n.Children {
		chStr += ch.String()
		if i != len(n.Children)-1 {
			chStr += " , "
		}
	}
	return fmt.Sprintf("N{<%s>%s; ch: [%s]}",
		nodeTypeStrings[n.Type], n.Value, chStr)
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

type EntityFilterDSLParser struct {
	lexer *EntityFilterDSLLexer
	token EntityFilterDSLToken
}

func (p *EntityFilterDSLParser) Parse(input string) (*Node, error) {
	p.lexer = &EntityFilterDSLLexer{}
	p.lexer.Init(strings.NewReader(input))
	p.token = p.lexer.Lex()
	node, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.token != EOF {
		return nil, errors.New("unexpected token after expression")
	}
	return node, nil
}

func (p *EntityFilterDSLParser) parseExpr() (*Node, error) {
	node := &Node{Type: NodeExpr}
	child, err := p.parsePredicateExpr()
	if err != nil {
		return nil, err
	}
	node.AddChild(child)

	if p.token == Semicolon {
		p.token = p.lexer.Lex()
		child, err := p.parseSortExpr()
		if err != nil {
			return nil, err
		}
		node.AddChild(child)
	}

	return node, nil
}

func (p *EntityFilterDSLParser) parsePredicateExpr() (*Node, error) {
	node := &Node{Type: NodePredicateExpr}

	if p.token == Not {
		node.AddChild(&Node{Type: NodeNot})
		p.token = p.lexer.Lex()
	}

	funcNode, err := p.parseFunction()
	if err != nil {
		return nil, err
	}
	node.AddChild(funcNode)

	if p.token == And || p.token == Or {
		op := &Node{Type: NodeType(p.token)}
		node.AddChild(op)
		p.token = p.lexer.Lex()
		child, err := p.parsePredicateExpr()
		if err != nil {
			return nil, err
		}
		op.AddChild(child)
	}

	return node, nil
}

func (p *EntityFilterDSLParser) parseFunction() (*Node, error) {
	if p.token != Function {
		return nil, fmt.Errorf("expected function, got: %v", p.token)
	}

	node := &Node{Type: NodeFunction, Value: p.lexer.TokenText()}
	p.token = p.lexer.Lex()

	if p.token != OpenParen {
		return nil, fmt.Errorf("expected open parenthesis, got: %v", p.token)
	}
	p.token = p.lexer.Lex()

	for p.token == Identifier {
		node.AddChild(&Node{Type: NodeIdentifier, Value: p.lexer.TokenText()})
		p.token = p.lexer.Lex()
		if p.token == Comma {
			p.token = p.lexer.Lex()
		}
	}

	if p.token != CloseParen {
		return nil, fmt.Errorf("expected close parenthesis, got: %v", p.token)
	}
	p.token = p.lexer.Lex()

	return node, nil
}

func (p *EntityFilterDSLParser) parseSortExpr() (*Node, error) {
	node := &Node{Type: NodeSortExpr}
	child, err := p.parseFunction()
	if err != nil {
		return nil, err
	}
	node.AddChild(child)
	return node, nil
}
