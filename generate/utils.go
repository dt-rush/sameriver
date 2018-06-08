package generate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func readSourceFile(srcFileName string) (
	*ast.File, []byte) {

	src, err := ioutil.ReadFile(srcFileName)
	if err != nil {
		panic(err)
	}
	astFile, err := parser.ParseFile(
		token.NewFileSet(), "", src, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	return astFile, src
}
