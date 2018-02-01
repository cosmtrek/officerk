package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Comm fields for all models
type Comm struct {
	ID        uint       `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// NullString ...
type NullString struct {
	sql.NullString
}

// MarshalJSON ...
func (v *NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

// UnmarshalJSON ...
func (v *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
	}
	return nil
}
