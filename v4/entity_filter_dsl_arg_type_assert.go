package sameriver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type DSLArgTypeAssertionFunc func(arg string, resolver IdentifierResolver) (any, error)

var typeAssertionFuncs = map[string]map[string]DSLArgTypeAssertionFunc{
	"IdentResolve": {
		"*Entity": func(arg string, resolver IdentifierResolver) (any, error) {
			entity, ok := resolver.Resolve(arg).(*Entity)
			if !ok {
				return nil, errors.New("type assert failed")
			}
			return entity, nil
		},
		"*Vec2D": func(arg string, resolver IdentifierResolver) (any, error) {
			vec2D, ok := resolver.Resolve(arg).(*Vec2D)
			if !ok {
				return nil, errors.New("type assert failed")
			}
			return vec2D, nil
		},
		"[]*Vec2D": func(arg string, resolver IdentifierResolver) (any, error) {
			vec2Ds, ok := resolver.Resolve(arg).([]*Vec2D)
			if !ok {
				return nil, errors.New("type assert failed")
			}
			return vec2Ds, nil
		},
		"*EventPredicate": func(arg string, resolver IdentifierResolver) (any, error) {
			eventPredicate, ok := resolver.Resolve(arg).(*EventPredicate)
			if !ok {
				return nil, errors.New("type assert failed")
			}
			return eventPredicate, nil
		},
		// Add more types here...
	},
	// Add other rules for other parametrized types here...
}

func ExtractTypesFromSignature(signature string) ([]string, error) {
	re := regexp.MustCompile(`\((.+)\)`)
	matches := re.FindStringSubmatch(signature)
	if len(matches) < 2 {
		return nil, fmt.Errorf("malformed signature: %s", signature)
	}
	typesStr := matches[1]
	types := strings.Split(typesStr, ", ")
	return types, nil
}

/*
DSLAssertArgTypes is responsible for checking if the arguments match the
expected types in a function signature. It uses the typeAssertionFuncs map to
perform type assertions for parametrized types. When a parametrized type is
encountered in the function signature, the type assertion function
corresponding to the type is called, and it returns the resolved value if the
assertion is successful.

a signature is a string like

WithinDistance(IdentResolve<*Entity>, int)
InPolygon(IdentResolve<[]*Vec2D>)
Closest(IdentResolve<*Entity>)
WithinRect(IdentResolve<*Vec2D>,IdentResolve<*Vec2D>)
*/
func DSLAssertArgTypes(signature string, args []string, resolver IdentifierResolver) ([]any, error) {
	expectedTypes, err := ExtractTypesFromSignature(signature)
	if err != nil {
		return nil, err
	}

	if len(args) != len(expectedTypes) {
		return nil, fmt.Errorf("number of arguments does not match number of expected types for function %s", signature)
	}

	resolved := make([]any, len(args))
	for i, arg := range args {
		parts := strings.Split(expectedTypes[i], "<")
		if typeAssertionFuncsMap, ok := typeAssertionFuncs[parts[0]]; ok && len(parts) > 1 {
			typeName := strings.TrimSuffix(parts[1], ">")
			if typeAssertionFunc, ok := typeAssertionFuncsMap[typeName]; ok {
				value, err := typeAssertionFunc(arg, resolver)
				if err != nil {
					return nil, fmt.Errorf("error for %s(%s): expected %s for argument %s, but %s", signature, strings.Join(args, ", "), expectedTypes[i], arg, err)
				}
				resolved[i] = value
			} else {
				return nil, fmt.Errorf("unsupported type in signature: %s", expectedTypes[i])
			}
		} else {
			switch expectedTypes[i] {
			case "string":
				resolved[i] = arg
			case "int":
				v, err := strconv.Atoi(arg)
				if err != nil {
					return nil, fmt.Errorf("expected int for argument %s, got: %s", arg, err)
				}
				resolved[i] = v
			case "float64":
				v, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					return nil, fmt.Errorf("expected float64 for argument %s, got: %s", arg, err)
				}
				resolved[i] = v
			case "[]string":
				resolved[i] = args
			default:
				return nil, fmt.Errorf("unsupported type in signature: %s", expectedTypes[i])
			}
		}
	}

	return resolved, nil
}

func DSLAssertOverloadedArgTypes(signatures []string, args []string, resolver IdentifierResolver) ([]interface{}, int, error) {
	for i, signature := range signatures {
		argsTyped, err := DSLAssertArgTypes(signature, args, resolver)
		if err == nil {
			return argsTyped, i, nil
		}
	}
	return nil, -1, fmt.Errorf("no matching signature found for overloaded method")
}
