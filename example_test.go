package uuid_test

import (
	"encoding/json"
	"fmt"

	guuid "github.com/google/uuid"
	"github.com/freddierice/uuid"
)

// User demonstrates how to use the ID type in a struct
type User struct {
	ID       uuid.ID         `json:"id" db:"id"`
	Name     string           `json:"name" db:"name"`
	Email    string           `json:"email" db:"email"`
	ManagerID uuid.NullableID `json:"manager_id" db:"manager_id"` // Optional foreign key
}

func ExampleID() {
	// Create a new ID
	id := uuid.New()
	fmt.Printf("UUID string: %s\n", id.String())
	fmt.Printf("Short string: %s\n", id.ShortString())
	
	// Create from existing UUID
	existingUUID := guuid.New()
	idFromUUID := uuid.FromUUID(existingUUID)
	fmt.Printf("From UUID: %s\n", idFromUUID.ShortString())
	
	// Parse from shortuuid string
	parsed, _ := uuid.Parse(id.ShortString())
	fmt.Printf("Parsed equals original: %t\n", parsed.Equal(id))
	
	// Output will vary due to random UUIDs, but structure will be:
	// UUID string: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	// Short string: xxxxxxxxxxxxxxxxxxxxxx
	// From UUID: xxxxxxxxxxxxxxxxxxxxxx
	// Parsed equals original: true
}

func ExampleID_json() {
	// Create a manager
	managerID := uuid.New()
	
	// Create a user with an ID and optional manager
	user := User{
		ID:        uuid.New(),
		Name:      "John Doe",
		Email:     "john@example.com",
		ManagerID: uuid.NewNullable(managerID), // Has a manager
	}
	
	// Create a user without a manager
	userNoManager := User{
		ID:    uuid.New(),
		Name:  "Jane Smith", 
		Email: "jane@example.com",
		// ManagerID will be invalid/null by default
	}
	
	// Marshal to JSON - IDs will be encoded as shortuuids, null manager_id as null
	jsonData, _ := json.Marshal(user)
	fmt.Printf("JSON with manager: %s\n", jsonData)
	
	jsonData2, _ := json.Marshal(userNoManager)
	fmt.Printf("JSON without manager: %s\n", jsonData2)
	
	// Unmarshal back - shortuuids will be decoded to UUIDs internally
	var parsed User
	json.Unmarshal(jsonData, &parsed)
	fmt.Printf("Same ID after JSON round-trip: %t\n", user.ID.Equal(parsed.ID))
	fmt.Printf("Manager ID preserved: %t\n", user.ManagerID.ID.Equal(parsed.ManagerID.ID))
	
	// Output will vary due to random IDs, but structure will be:
	// JSON with manager: {"id":"xxxxxxxxxxxxxxxxxxxxxx","name":"John Doe","email":"john@example.com","manager_id":"yyyyyyyyyyyyyyyyyyyyyy"}
	// JSON without manager: {"id":"xxxxxxxxxxxxxxxxxxxxxx","name":"Jane Smith","email":"jane@example.com","manager_id":null}
	// Same ID after JSON round-trip: true
	// Manager ID preserved: true
}
