package generate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func readSourceFile(srcFileName string) (
	*ast.File, []byte, error) {

	var fail = func(err error) (*ast.File, []byte, error) {
		return nil, []byte{}, err
	}

	src, err := ioutil.ReadFile(srcFileName)
	if err != nil {
		return fail(err)
	}
	astFile, err := parser.ParseFile(
		token.NewFileSet(), "", src, parser.AllErrors)
	if err != nil {
		return fail(err)
	}
	return astFile, src, nil
}
