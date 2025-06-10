package parser

import (
	"github.com/alecthomas/participle/v2"
	participleLexer "github.com/alecthomas/participle/v2/lexer"
)

type DSLFile struct {
	DatabaseDriver string   `"database" "driver" "=" @String`
	DatabaseURL    string   `"database" "url" "=" @String`
	Models         []*Model `@@*`
}

type Model struct {
	Model  string   `"model"`
	Name   string   `@Ident`
	LBrace string   `"{"`
	Fields []*Field `@@*`
	RBrace string   `"}"`
}

type Field struct {
	Name       string       `@Ident`
	Type       *Type        `@@`
	Directives []*Directive `@@*`
}

type Directive struct {
	Name string          `"@" @Ident`
	Args []*DirectiveArg `( "(" @@ ( "," @@ )* ")" )?`
}

type DirectiveArg struct {
	String *string  `  @String`
	Ident  *string  `| @Ident`
	Int    *int     `| @Int`
	Float  *float64 `| @Float`
}

type Type struct {
	Name    string `@Ident`
	IsArray bool   `(@"[" @"]")?`
}

// Configure the lexer to handle @ symbols and other tokens
var lexerRules = []participleLexer.SimpleRule{
	{Name: "Comment", Pattern: `//.*|/\*(.|\n)*?\*/`},
	{Name: "Whitespace", Pattern: `\s+`},
	{Name: "String", Pattern: `"[^"]*"`},
	{Name: "Float", Pattern: `[-+]?\d*\.\d+([eE][-+]?\d+)?`},
	{Name: "Int", Pattern: `[-+]?\d+`},
	{Name: "Ident", Pattern: `[a-zA-Z_]\w*`},
	{Name: "Punct", Pattern: `[@=(){}\[\],]`},
}

var stormLexer = participleLexer.MustSimple(lexerRules)

// Configure the parser with options to handle numeric values better
var Parser = participle.MustBuild[DSLFile](
	participle.Lexer(stormLexer),
	participle.Elide("Comment", "Whitespace"),
)
