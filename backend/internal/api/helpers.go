package api

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// uuidToString converts pgtype.UUID to string
func uuidToString(uuid pgtype.UUID) string {
	if !uuid.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}

// stringToPgUUID converts string to pgtype.UUID
func stringToPgUUID(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(s)
	return uuid, err
}

func ensureSlice[T any](items []T) []T {
	if items == nil {
		return []T{}
	}

	return items
}
