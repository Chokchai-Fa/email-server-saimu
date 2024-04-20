package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

func NewNullString(value string) NullString {
	return NullString{
		NullString: sql.NullString{
			String: value,
			Valid:  true,
		},
	}
}

// MarshalJSON for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// NullTime is an alias for mysql.NullTime data type
type NullTime struct {
	sql.NullTime
}

func NewNullTime(value time.Time) NullTime {
	return NullTime{
		NullTime: sql.NullTime{
			Time:  value,
			Valid: true,
		},
	}
}

// MarshalJSON for NullTime
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func NewNullInt64(value int64) NullInt64 {
	return NullInt64{
		NullInt64: sql.NullInt64{
			Int64: value,
			Valid: true,
		},
	}
}

// NullBool is an alias for sql.NullBool data type
type NullBool struct {
	sql.NullBool
}

// MarshalJSON for NullBool
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

func NewNullBool(value bool) NullBool {
	return NullBool{
		NullBool: sql.NullBool{
			Bool:  value,
			Valid: true,
		},
	}
}

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullFloat64
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

func NewNullFloat64(value float64) NullFloat64 {
	return NullFloat64{
		NullFloat64: sql.NullFloat64{
			Float64: value,
			Valid:   true,
		},
	}
}
