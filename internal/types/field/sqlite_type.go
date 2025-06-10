package field

func (ft FieldType) SQLiteType() string {
	switch ft.Kind {
	case KindInt:
		return "INTEGER"
	case KindFloat:
		return "REAL"
	case KindDecimal:
		return "NUMERIC"
	case KindBigInt:
		return "INTEGER"
	case KindString:
		return "TEXT"
	case KindText:
		return "TEXT"
	case KindChar:
		return "TEXT"
	case KindBoolean:
		return "INTEGER"
	case KindDateTime:
		return "TEXT"
	case KindDate:
		return "TEXT"
	case KindTime:
		return "TEXT"
	case KindTimestamp:
		return "TEXT"
	case KindBinary:
		return "BLOB"
	case KindJSON:
		return "TEXT"
	case KindUUID:
		return "TEXT"
	case KindCUID:
		return "TEXT"
	case KindPoint:
		return "TEXT"
	case KindCustom:
		return "TEXT"
	default:
		return "TEXT"
	}
}
