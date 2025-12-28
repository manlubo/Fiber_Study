package mapper

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Text

func TextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func TextValue(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// Timestamp

func TimePtr(t pgtype.Timestamp) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func TimeValue(t pgtype.Timestamp) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}
