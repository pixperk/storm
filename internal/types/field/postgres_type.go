package field

func (ft FieldType) PostgresType() string {
	switch ft.Kind {
	case KindInt:
		return "INTEGER"
	case KindFloat:
		return "DOUBLE PRECISION"
	case KindDecimal:
		if precision, scale := ft.GetPrecisionScale(); precision != "" && scale != "" {
			return "NUMERIC(" + precision + "," + scale + ")"
		}
		return "NUMERIC(10,2)"
	case KindBigInt:
		return "BIGINT"
	case KindString:
		length := returnLength(ft)
		return "VARCHAR(" + length + ")"
	case KindText:
		return "TEXT"
	case KindChar:
		length := returnLength(ft)
		return "CHAR(" + length + ")"
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
