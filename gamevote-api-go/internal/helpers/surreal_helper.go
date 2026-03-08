package helpers

import (
	"fmt"
	"strings"

	"github.com/surrealdb/surrealdb.go/pkg/models"
)

func ToRecordID(v any) (*models.RecordID, error) {
	switch val := v.(type) {
	case models.RecordID:
		return &val, nil
	case string:
		table, key, ok := strings.Cut(val, ":")
		if !ok {
			return nil, fmt.Errorf("invalid party id format: %s", val)
		}
		recordID := models.NewRecordID(table, key)
		return &recordID, nil
	default:
		return nil, fmt.Errorf("cannot convert to int")
	}
}
