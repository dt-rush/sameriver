package sameriver

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// these are really resolution funcs TODO rename
type DSLArgTypeAssertionFunc func(arg string, resolver IdentifierResolver) (any, error)

func AssertT[T any](x interface{}, expectedTypeStr string) (T, error) {
	t, ok := x.(T)
	if ok {
		return t, nil
	} else {
		return reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T), fmt.Errorf("type assertion failed: expected %s", expectedTypeStr)
	}
}

var IdentResolveTypeAssertMap = map[string]DSLArgTypeAssertionFunc{
	"*Entity": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[*Entity](resolver.Resolve(arg), "*Entity")
	},
	"string": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[string](resolver.Resolve(arg), "string")
	},
	"*Vec2D": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[*Vec2D](resolver.Resolve(arg), "*Vec2D")
	},
	"[]*Vec2D": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[[]*Vec2D](resolver.Resolve(arg), "[]*Vec2D")
	},
	"*EventPredicate": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[*EventPredicate](resolver.Resolve(arg), "*EventPredicate")
	},
	// Add more types here...
}

var typeResolveFuncs = map[string]map[string]DSLArgTypeAssertionFunc{
	"IdentResolve": IdentResolveTypeAssertMap,
	// Add other rules for other parametrized types than IdentResolve<T> here...
}

func ExtractTypesFromSignature(signature string) ([]string, error) {
	// Remove any whitespace before and after the commas
	signature = strings.TrimSpace(signature)

	// Make sure the signature is not empty
	if signature == "" {
		return nil, fmt.Errorf("malformed signature: %s", signature)
	}

	// Split the signature into individual types
	// we have to do it in this lexer-y way so that we can handle commas *inside* the
	// parametric types, like
	//
	// IdentResolve<*Entity>, TemporalFilter<*Entity,*Vec2D,*Vec2D>
	var types []string
	var currentType string
	var openBrackets int
	for i, r := range signature {
		if r == '<' {
			openBrackets++
		} else if r == '>' {
			openBrackets--
		} else if r == ',' && openBrackets == 0 {
			types = append(types, strings.TrimSpace(currentType))
			currentType = ""
			continue
		}
		currentType += string(r)
		if i == len(signature)-1 {
			types = append(types, strings.TrimSpace(currentType))
		}
	}

	return types, nil
}

/*
DSLAssertArgTypes is responsible for checking if the arguments match the
expected types in a function signature. It uses the typeResolveFuncs map to
perform type assertions for parametrized types. When a parametrized type is
encountered in the function signature, the type assertion function
corresponding to the type is called, and it returns the resolved value if the
assertion is successful.

a signature is a string like

IdentResolve<*Entity>, int
IdentResolve<[]*Vec2D>
IdentResolve<*Entity>
IdentResolve<*Vec2D>,IdentResolve<*Vec2D>
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
		// types are either simply defined in the switch statement below
		// or else are parametric in some way like IdentResolve<*Entity>
		parts := strings.Split(expectedTypes[i], "<")
		if typeResolveFuncsMap, ok := typeResolveFuncs[parts[0]]; ok && len(parts) > 1 {
			typeName := strings.TrimSuffix(parts[1], ">")
			if typeResolveFunc, ok := typeResolveFuncsMap[typeName]; ok {
				value, err := typeResolveFunc(arg, resolver)
				if err != nil {
					return nil, fmt.Errorf("error for %s(%s): expected %s for argument %s, but %s", signature, strings.Join(args, ", "), expectedTypes[i], arg, err)
				}
				resolved[i] = value
			} else {
				return nil, fmt.Errorf("unsupported type in signature: %s", expectedTypes[i])
			}
		} else {
			// regular types
			switch expectedTypes[i] {
			case "bool":
				v, err := strconv.ParseBool(arg)
				if err != nil {
					return nil, fmt.Errorf("expected bool for argument %s, got: %s", arg, err)
				}
				resolved[i] = v
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
			case "string":
				resolved[i] = arg
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
