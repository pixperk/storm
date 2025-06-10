package directive

type Directive struct {
	Kind DirectiveKind
	Args []string
}

func NewDirective(kind DirectiveKind, args []string) *Directive {
	return &Directive{
		Kind: kind,
		Args: args,
	}
}
