package field

type FieldKind int

const (
	// Numeric types
	KindInt     FieldKind = iota // INTEGER, INT, SMALLINT, TINYINT, MEDIUMINT, BIGINT
	KindFloat                    // FLOAT, DOUBLE, REAL, DECIMAL, NUMERIC
	KindDecimal                  // DECIMAL with fixed precision
	KindBigInt                   // For larger integers (64-bit)

	// String types
	KindString // VARCHAR, CHAR, TEXT
	KindText   // For larger text fields (TEXT, MEDIUMTEXT, LONGTEXT)
	KindChar   // Fixed-length character type

	// Boolean type
	KindBoolean // BOOLEAN, TINYINT(1)

	// Date and time types
	KindDateTime  // DATETIME
	KindDate      // DATE
	KindTime      // TIME
	KindTimestamp // TIMESTAMP

	// Binary data types
	KindBinary // BINARY, VARBINARY, BLOB

	// JSON and structured data
	KindJSON // JSON data type (supported in MySQL 5.7+, PostgreSQL, etc.)

	// Special types
	KindUUID  // UUID/GUID
	KindCUID  // CUID (Collision-resistant unique identifier)
	KindPoint // Geometric point (for GIS)

	// Custom types
	KindCustom // Custom type for user-defined field types
)

func (f FieldKind) String() string {
	switch f {
	case KindInt:
		return "Int"
	case KindFloat:
		return "Float"
	case KindDecimal:
		return "Decimal"
	case KindBigInt:
		return "BigInt"
	case KindString:
		return "String"
	case KindText:
		return "Text"
	case KindChar:
		return "Char"
	case KindBoolean:
		return "Boolean"
	case KindDateTime:
		return "DateTime"
	case KindDate:
		return "Date"
	case KindTime:
		return "Time"
	case KindTimestamp:
		return "Timestamp"
	case KindBinary:
		return "Binary"
	case KindJSON:
		return "JSON"
	case KindUUID:
		return "UUID"
	case KindCUID:
		return "CUID"
	case KindPoint:
		return "Point"
	case KindCustom:
		return "Custom"
	default:
		return "Unknown"
	}
}
