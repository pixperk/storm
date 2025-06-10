package field

func (k FieldKind) MySQLType() string {
	switch k {
	case KindInt:
		return "INT"
	case KindFloat:
		return "DOUBLE"
	case KindDecimal:
		return "DECIMAL(10,2)"
	case KindBigInt:
		return "BIGINT"
	case KindString:
		return "VARCHAR(255)"
	case KindText:
		return "TEXT"
	case KindChar:
		return "CHAR(255)"
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
		return "VARBINARY(255)"
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
