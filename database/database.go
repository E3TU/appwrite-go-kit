// Package database wraps the self-hosted Appwrite TablesDB surface used by Gin
// services.
//
// TablesDB is used rather than the legacy Databases API. The service is a thin
// layer over the SDK: rows are queried with the SDK query helpers.
package database

import (
	"errors"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/tablesdb"
)

// Service exposes the database operations the library supports.
type Service struct {
	sdk        *tablesdb.TablesDB
	databaseID string
}

// New builds a database Service bound to one database.
func New(c client.Client, databaseID string) *Service {
	return &Service{
		sdk:        tablesdb.New(c),
		databaseID: databaseID,
	}
}

// TablesDB returns the underlying SDK TablesDB service.
func (s *Service) TablesDB() *tablesdb.TablesDB {
	return s.sdk
}

// DatabaseID returns the database this service is bound to.
func (s *Service) DatabaseID() string {
	return s.databaseID
}

// ErrDisabled reports use of a database capability that is not configured.
var ErrDisabled = errors.New("database: capability disabled")
