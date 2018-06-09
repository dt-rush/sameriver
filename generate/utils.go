package generate

import (
	"fmt"
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

func getImportStringsFromFile(srcFileName string) (importStrings []string) {
	src, err := ioutil.ReadFile(srcFileName)
	if err != nil {
		panic(err)
	}
	return getImportStringsFromFileAsString(string(src))
}

func getImportStringsFromFileAsString(src string) (importStrings []string) {
	astFile, err := parser.ParseFile(
		token.NewFileSet(), "", src, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	for _, importSpec := range astFile.Imports {
		importStrings = append(importStrings,
			string(src[importSpec.Pos()-1:importSpec.End()-1]))
	}
	return importStrings
}

func importsBlockFromStrings(imports []string) string {
	importsBlock := "import (\n"
	for _, importToAdd := range imports {
		importsBlock += fmt.Sprintf("\t%s\n", importToAdd)
	}
	importsBlock += ")\n"
	return importsBlock
}
