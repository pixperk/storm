package directive

type DirectiveKind int

const (
	DirID DirectiveKind = iota
	DirAuto
	DirDefault
	DirUnique
	DirNullable
	DirHasMany
	DirBelongsTo
	DirHasOne
	DirIndex
	DirEnum
	DirUpdatedAt
	DirCreatedAt
	DirLength
	DirMin
	DirMax
	DirPrecision
	DirDefaultNow
	DirMap
	DirRelation
)

// String returns the string representation of the directive kind
func (d DirectiveKind) String() string {
	switch d {
	case DirID:
		return "id"
	case DirAuto:
		return "auto"
	case DirDefault:
		return "default"
	case DirUnique:
		return "unique"
	case DirNullable:
		return "nullable"
	case DirHasMany:
		return "hasmany"
	case DirBelongsTo:
		return "belongsto"
	case DirHasOne:
		return "hasone"
	case DirIndex:
		return "index"
	case DirEnum:
		return "enum"
	case DirUpdatedAt:
		return "updatedat"
	case DirCreatedAt:
		return "createdat"
	case DirLength:
		return "length"
	case DirMin:
		return "min"
	case DirMax:
		return "max"
	case DirPrecision:
		return "precision"
	case DirDefaultNow:
		return "defaultnow"
	case DirMap:
		return "map"
	case DirRelation:
		return "relation"
	default:
		return ""
	}
}
