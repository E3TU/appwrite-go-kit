package database

import "testing"

func TestDatabaseID(t *testing.T) {
	s := &Service{databaseID: "main"}
	if got := s.DatabaseID(); got != "main" {
		t.Fatalf("DatabaseID() = %q, want main", got)
	}
}
