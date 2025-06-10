package field

func (ft FieldType) MySQLType() string {
	switch ft.Kind {
	case KindInt:
		return "INT"
	case KindFloat:
		return "DOUBLE"
	case KindDecimal:
		if precision, scale := ft.GetPrecisionScale(); precision != "" && scale != "" {
			return "NUMERIC(" + precision + "," + scale + ")"
		}
		return "DECIMAL(10,2)"
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
		return "TINYINT(1)"
	case KindDateTime:
		return "DATETIME"
	case KindDate:
		return "DATE"
	case KindTime:
		return "TIME"
	case KindTimestamp:
		return "TIMESTAMP"
	case KindBinary:
		length := returnLength(ft)
		return "BINARY(" + length + ")"
	case KindJSON:
		return "JSON"
	case KindUUID:
		return "CHAR(36)"
	case KindCUID:
		return "CHAR(25)"
	case KindPoint:
		return "POINT"
	case KindCustom:
		return "VARCHAR(255)"
	default:
		return "VARCHAR(255)"
	}
}
