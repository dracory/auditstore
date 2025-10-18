package auditstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

func initDB(_ string) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		panic(err)
	}

	return db
}

func initStore(filepath string) (StoreInterface, error) {
	db := initDB(filepath)

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		AuditTableName:     "audit_log",
		AutomigrateEnabled: true,
		DebugEnabled:       false,
	})

	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("unexpected nil store")
	}

	return store, nil
}

func TestStoreAuditCreate(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	audit := NewRecord().
		SetObjectType("user").
		SetObjectID("user_123").
		SetAuthorID("admin_1")

	oldValue := map[string]interface{}{"name": "Old Name"}
	newValue := map[string]interface{}{"name": "New Name"}

	oldValueJSON, _ := json.Marshal(oldValue)
	newValueJSON, _ := json.Marshal(newValue)

	audit = audit.SetValueOld(string(oldValueJSON)).
		SetValueNew(string(newValueJSON))

	err = store.AuditCreate(audit)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if audit.ID() == "" {
		t.Error("audit ID should be set after creation")
	}
}

func TestStoreAuditGet(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Create test audit
	audit := NewRecord().
		SetObjectType("user").
		SetObjectID("user_123").
		SetAuthorID("admin_1")

	err = store.AuditCreate(audit)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Retrieve the audit
	foundAudit, err := store.AuditGet(audit.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if foundAudit == nil {
		t.Fatal("audit should be found")
	}

	if foundAudit.ID() != audit.ID() {
		t.Error("retrieved audit ID does not match")
	}

	if foundAudit.ObjectType() != "user" {
		t.Error("object type does not match")
	}
}

func TestStoreAuditList(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Create test audits
	for i := 0; i < 5; i++ {
		audit := NewRecord().
			SetObjectType("user").
			SetObjectID("user_123").
			SetAuthorID("admin_1")

		err = store.AuditCreate(audit)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	// Test listing all audits
	query := NewRecordQuery()
	audits, err := store.AuditList(query)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(audits) != 5 {
		t.Fatalf("expected 5 audits, got %d", len(audits))
	}

	// Test filtering by object type
	query = NewRecordQuery()
	query = query.SetObjectType("user")
	audits, err = store.AuditList(query)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(audits) != 5 {
		t.Fatalf("expected 5 user audits, got %d", len(audits))
	}

	// Test pagination
	query = NewRecordQuery()
	query = query.SetLimit(2).SetOffset(0)
	audits, err = store.AuditList(query)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(audits) != 2 {
		t.Fatalf("expected 2 audits, got %d", len(audits))
	}
}

func TestStoreAuditCount(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Create test audits
	for i := 0; i < 3; i++ {
		audit := NewRecord().
			SetObjectType("user").
			SetObjectID("user_123").
			SetAuthorID("admin_1")

		err = store.AuditCreate(audit)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	// Count all audits
	count, err := store.AuditCount(NewRecordQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 3 {
		t.Fatalf("expected 3 audits, got %d", count)
	}

	// Count with filter
	count, err = store.AuditCount(NewRecordQuery().SetObjectType("user"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 3 {
		t.Fatalf("expected 3 user audits, got %d", count)
	}
}

func TestStoreAuditDelete(t *testing.T) {
	store, err := initStore(":memory:")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Create test audit
	audit := NewRecord().
		SetObjectType("user").
		SetObjectID("user_123").
		SetAuthorID("admin_1")

	err = store.AuditCreate(audit)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Delete the audit
	err = store.AuditDelete(audit.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Try to retrieve deleted audit
	deletedAudit, err := store.AuditGet(audit.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if deletedAudit != nil {
		t.Error("audit should be deleted")
	}
}
