// These are util functions to turn generic types into pgx types.

package dbutil

import "github.com/jackc/pgx/v5/pgtype"

func String(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{
			Valid: false,
		}
	}

	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}

func Bool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{
			Valid: false,
		}
	}

	return pgtype.Bool{
		Bool:  *b,
		Valid: true,
	}
}

func Int32(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{
			Valid: false,
		}
	}

	return pgtype.Int4{
		Int32: *i,
		Valid: true,
	}
}
