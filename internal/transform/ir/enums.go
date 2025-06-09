package ir

type FieldType int

const (
	TypeInt FieldType = iota
	TypeString
	TypeBoolean
	TypeFloat
	TypeRelation // fallback for custom types (e.g. Post, User)
)

func (t FieldType) String() string {
	return [...]string{"Int", "String", "Boolean", "Float", "Relation"}[t]
}

type Directive int

const (
	DirID Directive = iota
	DirAuto
	DirUnique
	DirHasMany
	DirBelongsTo
)

func (d Directive) String() string {
	return [...]string{"id", "auto", "unique", "hasMany", "belongsTo"}[d]
}
