package auditstore

import "testing"

func TestNewRecordDefaults(t *testing.T) {
	record := NewRecord()

	if record == nil {
		t.Fatal("expected record instance")
	}

	if _, ok := record.(*recordImplementation); !ok {
		t.Fatalf("expected *recordImplementation, got %T", record)
	}

	if record.ID() == "" {
		t.Error("expected ID to be set")
	}

	if record.ObjectType() != "" {
		t.Errorf("expected empty object type, got %q", record.ObjectType())
	}

	if record.ObjectID() != "" {
		t.Errorf("expected empty object id, got %q", record.ObjectID())
	}

	if record.ValueOld() != "" {
		t.Errorf("expected empty old value, got %q", record.ValueOld())
	}

	if record.ValueNew() != "" {
		t.Errorf("expected empty new value, got %q", record.ValueNew())
	}

	if record.AuthorID() != "" {
		t.Errorf("expected empty author id, got %q", record.AuthorID())
	}

	if record.CreatedAt() == "" {
		t.Fatal("expected created_at to be set")
	}

	if created := record.CreatedAtCarbon(); created == nil {
		t.Fatal("expected created_at carbon instance")
	} else if created.ToDateTimeString() != record.CreatedAt() {
		t.Errorf("expected created_at carbon to match stored value, got %q", created.ToDateTimeString())
	}
}

func TestRecordSettersAndGetters(t *testing.T) {
	record := NewRecord()

	record.SetID("record-123")
	if record.ID() != "record-123" {
		t.Fatalf("expected ID to be 'record-123', got %q", record.ID())
	}

	returned := record.SetObjectType("user")
	if returned != record {
		t.Fatal("expected SetObjectType to return the same instance")
	}
	record = returned

	if next := record.SetObjectID("object-456"); next != record {
		t.Fatal("expected SetObjectID to return the same instance")
	} else {
		record = next
	}

	if next := record.SetValueOld("{\"name\":\"old\"}"); next != record {
		t.Fatal("expected SetValueOld to return the same instance")
	} else {
		record = next
	}

	if next := record.SetValueNew("{\"name\":\"new\"}"); next != record {
		t.Fatal("expected SetValueNew to return the same instance")
	} else {
		record = next
	}

	if next := record.SetAuthorID("author-789"); next != record {
		t.Fatal("expected SetAuthorID to return the same instance")
	} else {
		record = next
	}

	if next := record.SetCreatedAt("2025-01-02 03:04:05"); next != record {
		t.Fatal("expected SetCreatedAt to return the same instance")
	} else {
		record = next
	}

	if record.ObjectType() != "user" {
		t.Errorf("expected object type 'user', got %q", record.ObjectType())
	}

	if record.ObjectID() != "object-456" {
		t.Errorf("expected object id 'object-456', got %q", record.ObjectID())
	}

	if record.ValueOld() != "{\"name\":\"old\"}" {
		t.Errorf("expected old value '{\"name\":\"old\"}', got %q", record.ValueOld())
	}

	if record.ValueNew() != "{\"name\":\"new\"}" {
		t.Errorf("expected new value '{\"name\":\"new\"}', got %q", record.ValueNew())
	}

	if record.AuthorID() != "author-789" {
		t.Errorf("expected author id 'author-789', got %q", record.AuthorID())
	}

	if record.CreatedAt() != "2025-01-02 03:04:05" {
		t.Errorf("expected created_at '2025-01-02 03:04:05', got %q", record.CreatedAt())
	}

	if created := record.CreatedAtCarbon(); created == nil {
		t.Fatal("expected created_at carbon instance")
	} else if created.ToDateTimeString() != "2025-01-02 03:04:05" {
		t.Errorf("expected created_at carbon '2025-01-02 03:04:05', got %q", created.ToDateTimeString())
	}
}

func TestNewRecordFromExistingData(t *testing.T) {
	createdAt := "2025-06-15 10:11:12"
	data := map[string]string{
		COLUMN_ID:          "existing-id",
		COLUMN_OBJECT_TYPE: "order",
		COLUMN_OBJECT_ID:   "order-001",
		COLUMN_VALUE_OLD:   "old",
		COLUMN_VALUE_NEW:   "new",
		COLUMN_AUTHOR_ID:   "author",
		COLUMN_CREATED_AT:  createdAt,
	}

	record := NewRecordFromExistingData(data)

	if record.ID() != "existing-id" {
		t.Errorf("expected ID 'existing-id', got %q", record.ID())
	}

	if record.ObjectType() != "order" {
		t.Errorf("expected object type 'order', got %q", record.ObjectType())
	}

	if record.ObjectID() != "order-001" {
		t.Errorf("expected object id 'order-001', got %q", record.ObjectID())
	}

	if record.ValueOld() != "old" {
		t.Errorf("expected old value 'old', got %q", record.ValueOld())
	}

	if record.ValueNew() != "new" {
		t.Errorf("expected new value 'new', got %q", record.ValueNew())
	}

	if record.AuthorID() != "author" {
		t.Errorf("expected author id 'author', got %q", record.AuthorID())
	}

	if record.CreatedAt() != createdAt {
		t.Errorf("expected created_at '%s', got %q", createdAt, record.CreatedAt())
	}

	if created := record.CreatedAtCarbon(); created == nil {
		t.Fatal("expected created_at carbon instance")
	} else if created.ToDateTimeString() != createdAt {
		t.Errorf("expected created_at carbon '%s', got %q", createdAt, created.ToDateTimeString())
	}
}
