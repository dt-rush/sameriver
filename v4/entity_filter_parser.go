package sameriver

/*

import (
	"fmt"
	"strconv"
	"strings"
)

type EntityFilterDSLParser struct {
	lexer *EntityFilterDSLLexer
}

func (p *EntityFilterDSLParser) Parse(input string) (interface{}, error) {
	p.lexer = &EntityFilterDSLLexer{}
	p.lexer.Init(strings.NewReader(input))

	return p.parseExpression()
}

func (p *EntityFilterDSLParser) parseFunctionCall() (interface{}, error) {
	token := p.lexer.Lex()

	switch token {
	case HasTag:
		return p.parseHasTagFunctionCall()
	case CanBe:
		return p.parseCanBeFunctionCall()
	default:
		return nil, fmt.Errorf("unexpected token %s", p.lexer.TokenText())
	}
}

func (p *EntityFilterDSLParser) parseHasTagFunctionCall() (interface{}, error) {
	err := p.expect(OpenParen)
	if err != nil {
		return nil, err
	}

	token := p.lexer.Lex()
	if token != StringArgument {
		return nil, fmt.Errorf("expected string argument, but got %s", p.lexer.TokenText())
	}
	stringArg := p.lexer.TokenText()

	err = p.expect(CloseParen)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":   "HasTag",
		"tag":    stringArg,
		"negate": false,
	}, nil
}

func (p *EntityFilterDSLParser) parseCanBeFunctionCall() (interface{}, error) {
	err := p.expect(OpenParen)
	if err != nil {
		return nil, err
	}

	token := p.lexer.Lex()
	if token != StringArgument {
		return nil, fmt.Errorf("expected string argument, but got %s", p.lexer.TokenText())
	}
	stringArg := p.lexer.TokenText()

	err = p.expect(Comma)
	if err != nil {
		return nil, err
	}

	token = p.lexer.Lex()
	if token != IntegerArgument {
		return nil, fmt.Errorf("expected integer argument, but got %s", p.lexer.TokenText())
	}
	intArg, err := strconv.Atoi(p.lexer.TokenText())
	if err != nil {
		return nil, fmt.Errorf("expected integer argument, but got %s", p.lexer.TokenText())
	}

	err = p.expect(CloseParen)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":   "CanBe",
		"tag":    stringArg,
		"amount": intArg,
	}, nil
}

func (p *EntityFilterDSLParser) parseEntityReference() (interface{}, error) {
	token := p.lexer.Lex()

	switch token {
	case SelfReference:
		return "self", nil
	case MindReference:
		return p.parseMindReference()
	case BlackboardReference:
		return p.parseBlackboardReference()
	default:
		return nil, fmt.Errorf("unexpected token %s", p.lexer.TokenText())
	}
}

func (p *EntityFilterDSLParser) parseMindReference() (interface{}, error) {
	err := p.expect(Dot)
	if err != nil {
		return nil, err
	}

	token := p.lexer.Lex()
	if token != Identifier {
		return nil, fmt.Errorf("expected identifier, but got %s", p.lexer.TokenText())
	}
	identifier := p.lexer.TokenText()

	return map[string]interface{}{
		"type":       "MindReference",
		"identifier": identifier,
	}, nil
}

func (p *EntityFilterDSLParser) parseBlackboardReference() (interface{}, error) {
	err := p.expect(Identifier)
	if err != nil {
		return nil, err
	}
	blackboardName := p.lexer.TokenText()

	err = p.expect(Dot)
	if err != nil {
		return nil, err
	}

	token := p.lexer.Lex()
	if token != Identifier {
		return nil, fmt.Errorf("expected identifier, but got %s", p.lexer.TokenText())
	}
	identifier := p.lexer.TokenText()

	return map[string]interface{}{
		"type":           "BlackboardReference",
		"blackboardName": blackboardName,
		"identifier":     identifier,
	}, nil
}

func (p *EntityFilterDSLParser) expect(expectedToken EntityFilterDSLToken) error {
	token := p.lexer.Lex()
	if token != expectedToken {
		return fmt.Errorf("expected %s, but got %s", expectedToken, p.lexer.TokenText())
	}
	return nil
}

func (p *EntityFilterDSLParser) parse() (interface{}, error) {
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if !p.isEOF() {
		return nil, fmt.Errorf("unexpected token %s", p.lexer.TokenText())
	}

	return expression, nil
}

func (p *EntityFilterDSLParser) isEOF() bool {
	return p.lexer.Peek() == EOF
}

func (p *EntityFilterDSLParser) parseExpression() (interface{}, error) {
	token := p.lexer.Lex()

	switch token {
	case EntityFilterExpression:
		p.lexer.Unscan()
		return p.parseEntityFilterExpression()
	case OrderingExpression:
		p.lexer.Unscan()
		return p.parseOrderingExpression()
	default:
		return nil, fmt.Errorf("unexpected token %s", p.lexer.TokenText())
	}
}

func (p *EntityFilterDSLParser) parseEntityFilterExpression() (interface{}, error) {
	var filters []interface{}

	functionCall, err := p.parseFunctionCall()
	if err != nil {
		return nil, err
	}
	filters = append(filters, functionCall)

	for {
		token := p.lexer.Lex()
		if token != And {
			p.lexer.Unscan()
			break
		}

		functionCall, err := p.parseFunctionCall()
		if err != nil {
			return nil, err
		}
		filters = append(filters, functionCall)
	}

	return filters, nil
}

func (p *EntityFilterDSLParser) parseOrderingExpression() (interface{}, error) {
	err := p.expect(OpenParen)
	if err != nil {
		return nil, err
	}

	entityReference, err := p.parseEntityReference()
	if err != nil {
		return nil, err
	}

	err = p.expect(Comma)
	if err != nil {
		return nil, err
	}

	expression, err := p.parseEntityFilterExpression()
	if err != nil {
		return nil, err
	}

	err = p.expect(CloseParen)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":             "Closest",
		"entityReference":  entityReference,
		"orderingCriteria": expression,
	}, nil
}
*/
