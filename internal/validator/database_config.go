package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pixperk/storm/internal/transform/ir"
)

// validateDatabaseConfig validates the database configuration
func validateDatabaseConfig(irData *ir.IR) error {
	errList := new(multierror.Error)
	// Check if database driver is provided
	if irData.DatabaseDriver == "" {
		errList = multierror.Append(errList, errors.New("database driver is required"))
	} else {
		// Remove quotation marks from the driver name if present
		driverName := strings.Trim(irData.DatabaseDriver, "\"'")

		// Validate supported database drivers
		switch driverName {
		case "mysql", "postgres", "postgresql", "sqlite", "sqlite3":
			// Valid drivers, no error
		default:
			errList = multierror.Append(errList, fmt.Errorf("unsupported database driver: %s", driverName))
		}
	}

	// Check if database URL is provided
	if irData.DatabaseURL == "" {
		errList = multierror.Append(errList, errors.New("database URL is required"))
	} else {
		// Basic URL validation (could be more sophisticated)
		if !strings.Contains(irData.DatabaseURL, "://") {
			errList = multierror.Append(errList, fmt.Errorf("invalid database URL format: %s", irData.DatabaseURL))
		}
	}

	return errList.ErrorOrNil()
}
