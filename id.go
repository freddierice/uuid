package uuid

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	guuid "github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

// ID represents a UUID that can be marshaled as a shortuuid for JSON
// and as a regular UUID for database operations
type ID struct {
	guuid.UUID
}

// New creates a new random ID
func New() ID {
	return ID{UUID: guuid.New()}
}

// FromUUID creates an ID from an existing UUID
func FromUUID(u guuid.UUID) ID {
	return ID{UUID: u}
}

// Parse parses a shortuuid string into an ID
func Parse(s string) (ID, error) {
	u, err := shortuuid.DefaultEncoder.Decode(s)
	if err != nil {
		return ID{}, fmt.Errorf("failed to decode shortuuid: %w", err)
	}
	return ID{UUID: u}, nil
}

// FromString parses a standard UUID string into an ID
func FromString(s string) (ID, error) {
	u, err := guuid.Parse(s)
	if err != nil {
		return ID{}, fmt.Errorf("failed to parse uuid: %w", err)
	}
	return ID{UUID: u}, nil
}

// String returns the standard UUID string representation
func (id ID) String() string {
	return id.UUID.String()
}

// ShortString returns the shortuuid string representation
func (id ID) ShortString() string {
	return shortuuid.DefaultEncoder.Encode(id.UUID)
}

// MarshalJSON implements json.Marshaler to encode as shortuuid
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.ShortString())
}

// UnmarshalJSON implements json.Unmarshaler to decode from shortuuid
func (id *ID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	
	parsed, err := Parse(s)
	if err != nil {
		return err
	}
	
	*id = parsed
	return nil
}

// Value implements driver.Valuer for database storage as UUID
func (id ID) Value() (driver.Value, error) {
	return id.UUID.String(), nil
}

// Scan implements sql.Scanner for database retrieval from UUID
func (id *ID) Scan(value interface{}) error {
	if value == nil {
		*id = ID{}
		return nil
	}
	
	switch v := value.(type) {
	case string:
		u, err := guuid.Parse(v)
		if err != nil {
			return fmt.Errorf("failed to scan UUID string: %w", err)
		}
		*id = ID{UUID: u}
		return nil
	case []byte:
		u, err := guuid.Parse(string(v))
		if err != nil {
			return fmt.Errorf("failed to scan UUID bytes: %w", err)
		}
		*id = ID{UUID: u}
		return nil
	case guuid.UUID:
		*id = ID{UUID: v}
		return nil
	default:
		return fmt.Errorf("cannot scan %T into ID", value)
	}
}

// IsZero returns true if the ID is the zero UUID
func (id ID) IsZero() bool {
	return id.UUID == guuid.Nil
}

// Equal returns true if two IDs are equal
func (id ID) Equal(other ID) bool {
	return id.UUID == other.UUID
}

// NullableID represents an ID that can be null in the database
type NullableID struct {
	ID    ID
	Valid bool
}

// NewNullable creates a NullableID with a valid ID
func NewNullable(id ID) NullableID {
	return NullableID{ID: id, Valid: !id.IsZero()}
}

// NullableFromPtr creates a NullableID from a pointer
func NullableFromPtr(id *ID) NullableID {
	if id == nil || id.IsZero() {
		return NullableID{Valid: false}
	}
	return NullableID{ID: *id, Valid: true}
}

// Ptr returns a pointer to the ID if valid, nil otherwise
func (nid NullableID) Ptr() *ID {
	if !nid.Valid {
		return nil
	}
	return &nid.ID
}

// MarshalJSON implements json.Marshaler for NullableID
func (nid NullableID) MarshalJSON() ([]byte, error) {
	if !nid.Valid {
		return json.Marshal(nil)
	}
	return nid.ID.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler for NullableID
func (nid *NullableID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nid.Valid = false
		return nil
	}
	
	var id ID
	if err := id.UnmarshalJSON(data); err != nil {
		return err
	}
	
	nid.ID = id
	nid.Valid = true
	return nil
}

// Value implements driver.Valuer for database storage
func (nid NullableID) Value() (driver.Value, error) {
	if !nid.Valid {
		return nil, nil
	}
	return nid.ID.Value()
}

// Scan implements sql.Scanner for database retrieval
func (nid *NullableID) Scan(value interface{}) error {
	if value == nil {
		nid.Valid = false
		return nil
	}
	
	var id ID
	if err := id.Scan(value); err != nil {
		return err
	}
	
	nid.ID = id
	nid.Valid = true
	return nil
}
