package parser

import "github.com/alecthomas/participle/v2"

type DSLFile struct {
	Models []*Model `@@*`
}

type Model struct {
	Model  string   `"model"`
	Name   string   `@Ident`
	LBrace string   `"{"`
	Fields []*Field `@@*`
	RBrace string   `"}"`
}

type Field struct {
	Name       string   `@Ident`
	Type       *Type    `@@`
	Directives []string `("@" @Ident)*`
}

type Type struct {
	Name    string `@Ident`
	IsArray bool   `(@"[" @"]")?`
}

var Parser = participle.MustBuild[DSLFile]()
