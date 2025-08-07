package sqlerr

// Code describes a specific type of database error.
// The value Other is reported when an error does not map to any of the defined codes.
type Code string

const (
	// Other is reported when an error does not map to any of the defined codes.
	Other Code = "other"

	// NotNullViolation is reported when a not null constraint would be violated.
	NotNullViolation Code = "not_null_violation"

	// ForeignKeyViolation is reported when a foreign key constraint would be violated.
	ForeignKeyViolation Code = "foreign_key_violation"

	// UniqueViolation is reported when a unique constraint would be violated.
	UniqueViolation Code = "unique_violation"

	// CheckViolation is reported when a check constraint would be violated.
	CheckViolation Code = "check_violation"

	// ExcludeViolation is reported when an exclusion constraint would be violated.
	ExcludeViolation Code = "exclude_violation"

	// TransactionFailed is reported when running a command in a failed transaction,
	// due to some previous command failure.
	TransactionFailed Code = "transaction_failed"

	// DeadlockDetected is reported when a deadlock is detected.
	// Deadlock detection is done on a best-effort basis and not all deadlocks
	// can be detected.
	DeadlockDetected Code = "deadlock_detected"

	// TooManyConnections is reported when the database rejects a connection request
	// due to reaching the maximum number of connections.
	// This is different from blocking waiting on a connection pool.
	TooManyConnections Code = "too_many_connections"
)

// MapCode maps an underlying database error to a Code.
func MapCode(code string) Code {
	switch code {
	case "23502":
		return NotNullViolation
	case "23503":
		return ForeignKeyViolation
	case "23505":
		return UniqueViolation
	case "23514":
		return CheckViolation
	case "23P01":
		return ExcludeViolation
	case "25P02":
		return TransactionFailed
	case "40P01":
		return DeadlockDetected
	case "53300":
		return TooManyConnections
	default:
		return Other
	}
}

// Severity defines the severity of a database error.
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityFatal   Severity = "FATAL"
	SeverityPanic   Severity = "PANIC"
	SeverityWarning Severity = "WARNING"
	SeverityNotice  Severity = "NOTICE"
	SeverityDebug   Severity = "DEBUG"
	SeverityInfo    Severity = "INFO"
	SeverityLog     Severity = "LOG"
)

// MapSeverity maps the severity string from the underlying database
// to a Severity.
func MapSeverity(severity string) Severity {
	switch severity {
	case "ERROR":
		return SeverityError
	case "FATAL":
		return SeverityFatal
	case "PANIC":
		return SeverityPanic
	case "WARNING":
		return SeverityWarning
	case "NOTICE":
		return SeverityNotice
	case "DEBUG":
		return SeverityDebug
	case "INFO":
		return SeverityInfo
	case "LOG":
		return SeverityLog
	default:
		return SeverityError
	}
}

// Error represents an error reported by the database server.
// It's not guaranteed all errors reported by database functions will be of this type;
// it is only returned when the database reports an error.
type Error struct {
	// Code defines the general class of the error.
	Code Code

	// Severity is the severity of the error.
	Severity Severity

	// DatabaseCode is the database server-specific error code.
	DatabaseCode string

	// Message: the primary human-readable error message.
	Message string

	// SchemaName: if the error was associated with a specific database object,
	// the name of the schema containing that object, if any.
	SchemaName string

	// TableName: if the error was associated with a specific table, the name of the table.
	TableName string

	// ColumnName: if the error was associated with a specific table column,
	// the name of the column.
	ColumnName string

	// DataTypeName: if the error was associated with a specific data type,
	// the name of the data type.
	DataTypeName string

	// ConstraintName: if the error was associated with a specific constraint,
	// the name of the constraint.
	ConstraintName string

	// driverErr is the underlying error from the driver.
	driverErr error
}

func (pe *Error) Error() string {
	return string(pe.Severity) + ": " + pe.Message + " (Code " + string(pe.Code) + ": SQLSTATE " + pe.DatabaseCode + ")"
}

func (pe *Error) Unwrap() error {
	return pe.driverErr
}
