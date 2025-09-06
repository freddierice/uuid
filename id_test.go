package uuid

import (
	"encoding/json"
	"testing"

	guuid "github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

func TestNew(t *testing.T) {
	id := New()
	if id.IsZero() {
		t.Error("New should not return zero UUID")
	}
}

func TestFromUUID(t *testing.T) {
	u := guuid.New()
	id := FromUUID(u)
	if id.UUID != u {
		t.Error("FromUUID should preserve UUID")
	}
}

func TestParse(t *testing.T) {
	// Create a known UUID and its shortuuid representation
	originalUUID := guuid.New()
	shortStr := shortuuid.DefaultEncoder.Encode(originalUUID)
	
	// Parse the shortuuid back to ID
	parsed, err := Parse(shortStr)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if parsed.UUID != originalUUID {
		t.Error("Parse should recover original UUID")
	}
}

func TestFromString(t *testing.T) {
	uuidStr := "550e8400-e29b-41d4-a716-446655440000"
	id, err := FromString(uuidStr)
	if err != nil {
		t.Fatalf("FromString failed: %v", err)
	}
	
	if id.String() != uuidStr {
		t.Errorf("Expected %s, got %s", uuidStr, id.String())
	}
}

func TestString(t *testing.T) {
	u := guuid.New()
	id := FromUUID(u)
	
	if id.String() != u.String() {
		t.Error("String() should return UUID string representation")
	}
}

func TestShortString(t *testing.T) {
	u := guuid.New()
	id := FromUUID(u)
	expected := shortuuid.DefaultEncoder.Encode(u)
	
	if id.ShortString() != expected {
		t.Error("ShortString() should return shortuuid representation")
	}
}

func TestJSONMarshaling(t *testing.T) {
	id := New()
	
	// Marshal to JSON
	data, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	
	// Should be a shortuuid string
	var shortStr string
	err = json.Unmarshal(data, &shortStr)
	if err != nil {
		t.Fatalf("Failed to unmarshal as string: %v", err)
	}
	
	if shortStr != id.ShortString() {
		t.Error("JSON should contain shortuuid representation")
	}
}

func TestJSONUnmarshaling(t *testing.T) {
	originalID := New()
	shortStr := originalID.ShortString()
	
	// Create JSON with shortuuid
	data, _ := json.Marshal(shortStr)
	
	// Unmarshal back to ID
	var id ID
	err := json.Unmarshal(data, &id)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}
	
	if !id.Equal(originalID) {
		t.Error("JSON round-trip should preserve ID")
	}
}

func TestJSONRoundTrip(t *testing.T) {
	original := New()
	
	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	
	// Unmarshal back
	var parsed ID
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	
	if !original.Equal(parsed) {
		t.Error("JSON round-trip should preserve ID")
	}
}

func TestIsZero(t *testing.T) {
	var zero ID
	if !zero.IsZero() {
		t.Error("Zero ID should return true for IsZero()")
	}
	
	nonZero := New()
	if nonZero.IsZero() {
		t.Error("Non-zero ID should return false for IsZero()")
	}
}

func TestEqual(t *testing.T) {
	id1 := New()
	id2 := FromUUID(id1.UUID)
	id3 := New()
	
	if !id1.Equal(id2) {
		t.Error("IDs with same UUID should be equal")
	}
	
	if id1.Equal(id3) {
		t.Error("IDs with different UUIDs should not be equal")
	}
}

func TestScanValue(t *testing.T) {
	original := New()
	
	// Get database value
	value, err := original.Value()
	if err != nil {
		t.Fatalf("Value() failed: %v", err)
	}
	
	// Scan it back
	var scanned ID
	err = scanned.Scan(value)
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}
	
	if !original.Equal(scanned) {
		t.Error("Database round-trip should preserve ID")
	}
}

func TestScanNil(t *testing.T) {
	var id ID
	err := id.Scan(nil)
	if err != nil {
		t.Fatalf("Scan(nil) failed: %v", err)
	}
	
	if !id.IsZero() {
		t.Error("Scanning nil should result in zero ID")
	}
}

func TestNewNullable(t *testing.T) {
	id := New()
	nid := NewNullable(id)
	
	if !nid.Valid {
		t.Error("NullableID should be valid for non-zero ID")
	}
	
	if !nid.ID.Equal(id) {
		t.Error("NullableID should contain the original ID")
	}
	
	// Test with zero ID
	var zeroID ID
	nid2 := NewNullable(zeroID)
	if nid2.Valid {
		t.Error("NullableID should not be valid for zero ID")
	}
}

func TestNullableFromPtr(t *testing.T) {
	id := New()
	
	// Test with valid pointer
	nid := NullableFromPtr(&id)
	if !nid.Valid {
		t.Error("NullableID should be valid for non-nil pointer")
	}
	
	if !nid.ID.Equal(id) {
		t.Error("NullableID should contain the original ID")
	}
	
	// Test with nil pointer
	nid2 := NullableFromPtr(nil)
	if nid2.Valid {
		t.Error("NullableID should not be valid for nil pointer")
	}
	
	// Test with zero ID pointer
	var zeroID ID
	nid3 := NullableFromPtr(&zeroID)
	if nid3.Valid {
		t.Error("NullableID should not be valid for zero ID pointer")
	}
}

func TestNullableIDPtr(t *testing.T) {
	id := New()
	nid := NewNullable(id)
	
	ptr := nid.Ptr()
	if ptr == nil {
		t.Fatal("Ptr() should not return nil for valid NullableID")
	}
	
	if !ptr.Equal(id) {
		t.Error("Ptr() should return pointer to original ID")
	}
	
	// Test invalid NullableID
	var invalid NullableID
	ptr2 := invalid.Ptr()
	if ptr2 != nil {
		t.Error("Ptr() should return nil for invalid NullableID")
	}
}

func TestNullableIDJSONMarshaling(t *testing.T) {
	id := New()
	nid := NewNullable(id)
	
	// Test valid NullableID
	data, err := json.Marshal(nid)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	
	// Should be the same as marshaling the ID directly
	expectedData, _ := json.Marshal(id)
	if string(data) != string(expectedData) {
		t.Error("Valid NullableID should marshal same as ID")
	}
	
	// Test invalid NullableID
	var invalid NullableID
	data2, err := json.Marshal(invalid)
	if err != nil {
		t.Fatalf("JSON marshal of invalid failed: %v", err)
	}
	
	if string(data2) != "null" {
		t.Error("Invalid NullableID should marshal to null")
	}
}

func TestNullableIDJSONUnmarshaling(t *testing.T) {
	id := New()
	
	// Test unmarshaling valid ID
	data, _ := json.Marshal(id.ShortString())
	var nid NullableID
	err := json.Unmarshal(data, &nid)
	if err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}
	
	if !nid.Valid {
		t.Error("NullableID should be valid after unmarshaling valid ID")
	}
	
	if !nid.ID.Equal(id) {
		t.Error("NullableID should contain original ID after unmarshaling")
	}
	
	// Test unmarshaling null
	var nid2 NullableID
	err = json.Unmarshal([]byte("null"), &nid2)
	if err != nil {
		t.Fatalf("JSON unmarshal of null failed: %v", err)
	}
	
	if nid2.Valid {
		t.Error("NullableID should not be valid after unmarshaling null")
	}
}

func TestNullableIDDatabaseOperations(t *testing.T) {
	id := New()
	nid := NewNullable(id)
	
	// Test Value() for valid NullableID
	value, err := nid.Value()
	if err != nil {
		t.Fatalf("Value() failed: %v", err)
	}
	
	if value == nil {
		t.Error("Value() should not return nil for valid NullableID")
	}
	
	// Test Scan() for valid value
	var scanned NullableID
	err = scanned.Scan(value)
	if err != nil {
		t.Fatalf("Scan() failed: %v", err)
	}
	
	if !scanned.Valid {
		t.Error("Scanned NullableID should be valid")
	}
	
	if !scanned.ID.Equal(id) {
		t.Error("Database round-trip should preserve ID")
	}
	
	// Test invalid NullableID
	var invalid NullableID
	value2, err := invalid.Value()
	if err != nil {
		t.Fatalf("Value() failed for invalid: %v", err)
	}
	
	if value2 != nil {
		t.Error("Value() should return nil for invalid NullableID")
	}
	
	// Test scanning nil
	var scanned2 NullableID
	err = scanned2.Scan(nil)
	if err != nil {
		t.Fatalf("Scan(nil) failed: %v", err)
	}
	
	if scanned2.Valid {
		t.Error("Scanning nil should result in invalid NullableID")
	}
}

func TestRealUUIDPairs(t *testing.T) {
	testPairs := []struct {
		uuidStr   string
		shortUUID string
	}{
		{"8eca1fe1-a833-4de0-b487-b54785cc656e", "TR7kT7YgYjHN5ET8Yagssh"},
		{"dbf0e6e9-412d-4789-8042-b35d88d5fe80", "h9YNy27x6xdsqHTjMouA6s"},
		{"47934534-601a-4ea2-9f35-28d761dbe417", "EjtEAyHfuns4hG3SwZKMhc"},
		{"10fd962b-efd3-4725-b556-9a7caad7a40e", "53Kg5b2NMsBvKw8net9Mss"},
		{"00000000-0000-0000-0000-000000000000", "2222222222222222222222"},
	}

	for _, pair := range testPairs {
		t.Run(pair.uuidStr, func(t *testing.T) {
			// Test UUID string -> ID -> shortuuid string
			id, err := FromString(pair.uuidStr)
			if err != nil {
				t.Fatalf("FromString failed for %s: %v", pair.uuidStr, err)
			}
			
			shortStr := id.ShortString()
			if shortStr != pair.shortUUID {
				t.Errorf("Expected short UUID %s, got %s", pair.shortUUID, shortStr)
			}
			
			// Test shortuuid string -> ID -> UUID string
			parsedID, err := Parse(pair.shortUUID)
			if err != nil {
				t.Fatalf("Parse failed for %s: %v", pair.shortUUID, err)
			}
			
			uuidStr := parsedID.String()
			if uuidStr != pair.uuidStr {
				t.Errorf("Expected UUID %s, got %s", pair.uuidStr, uuidStr)
			}
			
			// Verify both IDs are equal
			if !id.Equal(parsedID) {
				t.Error("IDs from UUID and shortuuid should be equal")
			}
		})
	}
}

func TestRealUUIDPairsJSON(t *testing.T) {
	testPairs := []struct {
		uuidStr   string
		shortUUID string
	}{
		{"8eca1fe1-a833-4de0-b487-b54785cc656e", "TR7kT7YgYjHN5ET8Yagssh"},
		{"dbf0e6e9-412d-4789-8042-b35d88d5fe80", "h9YNy27x6xdsqHTjMouA6s"},
		{"47934534-601a-4ea2-9f35-28d761dbe417", "EjtEAyHfuns4hG3SwZKMhc"},
		{"10fd962b-efd3-4725-b556-9a7caad7a40e", "53Kg5b2NMsBvKw8net9Mss"},
		{"00000000-0000-0000-0000-000000000000", "2222222222222222222222"},
	}

	for _, pair := range testPairs {
		t.Run(pair.uuidStr, func(t *testing.T) {
			id, err := FromString(pair.uuidStr)
			if err != nil {
				t.Fatalf("FromString failed: %v", err)
			}
			
			// Test JSON marshaling - should produce shortuuid
			jsonData, err := json.Marshal(id)
			if err != nil {
				t.Fatalf("JSON marshal failed: %v", err)
			}
			
			var marshaledShort string
			err = json.Unmarshal(jsonData, &marshaledShort)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON as string: %v", err)
			}
			
			if marshaledShort != pair.shortUUID {
				t.Errorf("Expected marshaled JSON to contain %s, got %s", pair.shortUUID, marshaledShort)
			}
			
			// Test JSON unmarshaling - should decode shortuuid back to UUID
			var unmarshaledID ID
			err = json.Unmarshal(jsonData, &unmarshaledID)
			if err != nil {
				t.Fatalf("JSON unmarshal failed: %v", err)
			}
			
			if !id.Equal(unmarshaledID) {
				t.Error("JSON round-trip should preserve ID")
			}
			
			if unmarshaledID.String() != pair.uuidStr {
				t.Errorf("Expected UUID %s after unmarshaling, got %s", pair.uuidStr, unmarshaledID.String())
			}
		})
	}
}

func TestRealUUIDPairsDatabase(t *testing.T) {
	testPairs := []struct {
		uuidStr   string
		shortUUID string
	}{
		{"8eca1fe1-a833-4de0-b487-b54785cc656e", "TR7kT7YgYjHN5ET8Yagssh"},
		{"dbf0e6e9-412d-4789-8042-b35d88d5fe80", "h9YNy27x6xdsqHTjMouA6s"},
		{"47934534-601a-4ea2-9f35-28d761dbe417", "EjtEAyHfuns4hG3SwZKMhc"},
		{"10fd962b-efd3-4725-b556-9a7caad7a40e", "53Kg5b2NMsBvKw8net9Mss"},
		{"00000000-0000-0000-0000-000000000000", "2222222222222222222222"},
	}

	for _, pair := range testPairs {
		t.Run(pair.uuidStr, func(t *testing.T) {
			id, err := FromString(pair.uuidStr)
			if err != nil {
				t.Fatalf("FromString failed: %v", err)
			}
			
			// Test database Value() - should produce UUID string
			value, err := id.Value()
			if err != nil {
				t.Fatalf("Value() failed: %v", err)
			}
			
			valueStr, ok := value.(string)
			if !ok {
				t.Fatalf("Expected Value() to return string, got %T", value)
			}
			
			if valueStr != pair.uuidStr {
				t.Errorf("Expected Value() to return %s, got %s", pair.uuidStr, valueStr)
			}
			
			// Test database Scan() - should parse UUID string back to ID
			var scannedID ID
			err = scannedID.Scan(value)
			if err != nil {
				t.Fatalf("Scan() failed: %v", err)
			}
			
			if !id.Equal(scannedID) {
				t.Error("Database round-trip should preserve ID")
			}
			
			if scannedID.String() != pair.uuidStr {
				t.Errorf("Expected UUID %s after scanning, got %s", pair.uuidStr, scannedID.String())
			}
			
			if scannedID.ShortString() != pair.shortUUID {
				t.Errorf("Expected short UUID %s after scanning, got %s", pair.shortUUID, scannedID.ShortString())
			}
			
			// Test scanning as []byte (common database driver behavior)
			var scannedID2 ID
			err = scannedID2.Scan([]byte(pair.uuidStr))
			if err != nil {
				t.Fatalf("Scan([]byte) failed: %v", err)
			}
			
			if !id.Equal(scannedID2) {
				t.Error("Database []byte round-trip should preserve ID")
			}
		})
	}
}
