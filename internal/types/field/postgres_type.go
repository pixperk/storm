package field

func (k FieldKind) PostgresType() string {
	switch k {
	case KindInt:
		return "INTEGER"
	case KindFloat:
		return "DOUBLE PRECISION"
	case KindDecimal:
		return "NUMERIC(10,2)"
	case KindBigInt:
		return "BIGINT"
	case KindString:
		return "VARCHAR(255)"
	case KindText:
		return "TEXT"
	case KindChar:
		return "CHAR(255)"
	case KindBoolean:
		return "BOOLEAN"
	case KindDateTime:
		return "TIMESTAMP"
	case KindDate:
		return "DATE"
	case KindTime:
		return "TIME"
	case KindTimestamp:
		return "TIMESTAMP"
	case KindBinary:
		return "BYTEA"
	case KindJSON:
		return "JSONB"
	case KindUUID:
		return "UUID"
	case KindCUID:
		return "VARCHAR(25)"
	case KindPoint:
		return "POINT"
	case KindCustom:
		return "VARCHAR(255)"
	default:
		return "VARCHAR(255)"
	}
}
