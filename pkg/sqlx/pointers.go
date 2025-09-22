package sqlx

import "database/sql"

func Bool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{
			Valid: false,
		}
	}

	return sql.NullBool{
		Bool:  *b,
		Valid: true,
	}
}

func Int32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{
			Valid: false,
		}
	}

	return sql.NullInt32{
		Int32: *i,
		Valid: true,
	}
}
