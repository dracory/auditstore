# Auditstore

Auditstore is a lightweight Go module for persisting domain audit records with a relational database. It wraps a standard `*sql.DB` connection, provides a fluent `RecordInterface` model, and offers query helpers and automigration support so applications can track create/update/delete history with minimal boilerplate.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Key Concepts](#key-concepts)
- [API Reference](#api-reference)
- [Querying Records](#querying-records)
- [Schema Details](#schema-details)
- [Debugging & Logging](#debugging--logging)
- [Testing](#testing)
- [Development Notes](#development-notes)
- [Dependencies](#dependencies)

## Overview
Auditstore centers on the `storeImplementation` in `store_implementation.go`, which implements `StoreInterface` for creating, reading, listing, counting, and deleting audits. Records are modeled through `recordImplementation` in `record_implementation.go`, while query construction is handled by `RecordQuery` in `record_query.go`.

## Features
- **Automated schema management** via `AutoMigrate()` and SQL builders in `sqls.go` using `github.com/dracory/sb`.
- **Fluent record modeling** with `RecordInterface` setters that wrap a `dataobject.DataObject`.
- **Dialect-aware SQL generation** leveraging `github.com/doug-martin/goqu/v9` for inserts, selects, counts, and deletes.
- **Query builder** supporting filtering, ordering, pagination, and date ranges through `RecordQueryInterface`.
- **Debug logging** with structured output powered by `log/slog` when `EnableDebugMode(true)` is invoked.

## Installation
```bash
go get github.com/dracory/auditstore
```
Ensure your project uses Go `1.24.5` or newer (specified in `go.mod`).

## Quick Start
```go
package main

import (
    "database/sql"
    "log"

    auditstore "github.com/dracory/auditstore"
    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "file:audit.db?cache=shared&mode=rwc&parseTime=true")
    if err != nil {
        log.Fatal(err)
    }

    store, err := auditstore.NewStore(auditstore.NewStoreOptions{
        DB:                 db,
        AuditTableName:     "audit_log",
        AutomigrateEnabled: true,
        DebugEnabled:       false,
    })
    if err != nil {
        log.Fatal(err)
    }

    record := auditstore.NewRecord().
        SetObjectType("user").
        SetObjectID("user_123").
        SetAuthorID("admin_1").
        SetValueOld(`{"name":"Old Name"}`).
        SetValueNew(`{"name":"New Name"}`)

    if err := store.AuditCreate(record); err != nil {
        log.Fatal(err)
    }

    fetched, err := store.AuditGet(record.ID())
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Fetched audit %s created at %s", fetched.ID(), fetched.CreatedAt())
}
```

## Key Concepts
- **Store**: `NewStore(NewStoreOptions)` returns a `StoreInterface` backed by `storeImplementation`. It requires an initialized `*sql.DB` and the target table name.
- **Record**: `NewRecord()` creates a mutable audit entry. IDs and timestamps default to `uid.MicroUid()` and the current UTC time (via `github.com/dromara/carbon/v2`).
- **Query**: `NewRecordQuery()` produces a builder you can chain to filter, sort, and paginate audit results.

## API Reference
### StoreInterface (`store_interface.go`)
- **EnableDebugMode(debug bool)**: Toggle structured logging; `DebugEnable` aliases this for backward compatibility.
- **AuditCreate(record RecordInterface) error**: Inserts a new audit row, auto-populating `id` and `created_at` when missing.
- **AuditGet(id string) (RecordInterface, error)**: Retrieves a single record or `nil` if not found.
- **AuditList(query RecordQueryInterface) ([]RecordInterface, error)**: Returns all records matching the query.
- **AuditCount(query RecordQueryInterface) (int64, error)**: Counts rows for the given query.
- **AuditDelete(id string) error**: Removes the matching row.
- **AutoMigrate() error**: Creates the audit table using `sqlAuditTableCreate()` when absent.

### NewStoreOptions (`store_implementation.go`)
- **DB**: Required `*sql.DB` connection. Its driver name determines SQL dialect via `database.DatabaseType`.
- **AuditTableName**: Required table name (e.g., `audit_log`).
- **AutomigrateEnabled**: Auto-run `AutoMigrate()` during construction.
- **DebugEnabled**: Enable structured debug logging immediately.

### RecordInterface (`record_interface.go`)
- **SetObjectType(string)** / `ObjectType()`
- **SetObjectID(string)** / `ObjectID()`
- **SetValueOld(string)** / `ValueOld()`
- **SetValueNew(string)** / `ValueNew()`
- **SetAuthorID(string)** / `AuthorID()`
- **SetCreatedAt(string)** / `CreatedAt()`
- **CreatedAtCarbon()**: Returns `*carbon.Carbon` for time arithmetic.

## Querying Records
`RecordQueryInterface` (`record_query_interface.go`) exposes fluent setters. The concrete implementation in `record_query.go` validates inputs and converts them into a `goqu.SelectDataset` using `ToSelectDataset(driver, table)`.

```go
query := auditstore.NewRecordQuery().
    SetObjectType("user").
    SetAuthorID("admin_1").
    SetCreatedAfter(time.Now().Add(-24 * time.Hour)).
    SetOrderBy(auditstore.COLUMN_CREATED_AT, false).
    SetLimit(50).
    SetOffset(0)

records, err := store.AuditList(query)
if err != nil {
    // handle error
}

count, err := store.AuditCount(query)
```

Validation occurs via `RecordQuery.Validate()`, ensuring limits are non-negative and date ranges are consistent before SQL is generated.

## Schema Details
Column names are defined in `constants.go`:
- **`COLUMN_ID`** (`id`)
- **`COLUMN_OBJECT_TYPE`** (`object_type`)
- **`COLUMN_OBJECT_ID`** (`object_id`)
- **`COLUMN_VALUE_OLD`** (`value_old`)
- **`COLUMN_VALUE_NEW`** (`value_new`)
- **`COLUMN_AUTHOR_ID`** (`author_id`)
- **`COLUMN_CREATED_AT`** (`created_at`)

`sqls.go` uses `github.com/dracory/sb` to declare the table with appropriate column lengths and types. Indexes can be added by extending the builder section that currently contains commented placeholders.

## Debugging & Logging
- **Enable debug mode**: `store.EnableDebugMode(true)` switches the logger to debug level using `slog.NewTextHandler` with `LevelDebug`.
- **SQL output**: When debug is on, SQL statements and parameters are printed before execution across CRUD methods.

## Testing
`store_test.go` contains integration-style tests using `modernc.org/sqlite` in-memory databases:
- **`TestStoreAuditCreate`**: Ensures IDs are assigned and records persist.
- **`TestStoreAuditGet`**, **`TestStoreAuditList`**, **`TestStoreAuditCount`**, **`TestStoreAuditDelete`** cover the primary Store operations.
Run the suite with:
```bash
go test ./...
```

## Development Notes
- Ensure the target database driver is registered with `database/sql` (e.g., import `_ "modernc.org/sqlite"`).
- `AuditCreate` uses `database.Execute` from `github.com/dracory/database`, which expects context-aware execution; the package wraps `context.Background()` for convenience.
- Timestamps default to UTC via `carbon.Now(carbon.UTC)` inside `AuditCreate`, guaranteeing consistent storage across drivers.

## Dependencies
Auditstore relies on:
- **`github.com/doug-martin/goqu/v9`** for SQL construction.
- **`github.com/dracory/database`** for database helpers (`SelectToMapString`, `Execute`).
- **`github.com/dracory/dataobject`** for the record data container.
- **`github.com/dracory/uid`** for unique ID generation.
- **`github.com/dracory/sb`** for schema building.
- **`github.com/dromara/carbon/v2`** for time handling.
- **`modernc.org/sqlite`** in tests as an embedded driver; replace or supplement with your chosen driver in production.
