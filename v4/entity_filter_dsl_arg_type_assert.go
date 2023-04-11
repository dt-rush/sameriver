package sameriver

import (
	"fmt"
	"reflect"
	"regexp"
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
expected types in a function signature. It uses the typeResolveFuncs map to
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
