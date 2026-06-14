# Audit Store

[![Tests Status](https://github.com/dracory/auditstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/auditstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/auditstore)](https://goreportcard.com/report/github.com/dracory/auditstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/auditstore)](https://pkg.go.dev/github.com/dracory/auditstore)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/v/release/dracory/auditstore?include_prereleases&style=flat-square)](https://github.com/dracory/auditstore/releases)
[![Go version](https://img.shields.io/github/go-mod/go-version/dracory/auditstore)](https://github.com/dracory/auditstore)
[![codecov](https://codecov.io/gh/dracory/auditstore/branch/main/graph/badge.svg)](https://codecov.io/gh/dracory/auditstore)

Audit store keeps a clear story of every change that happens in your product. Instead of guessing who updated a record or when something broke, you get a searchable timeline showing the actor, the item they touched, and what changed. Plug it into your existing SQL database, drop audit entries wherever you mutate data, and use the built-in queries to review activity or restore confidence when issues appear.

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
Auditstore helps your application answer the question *"what changed and when?"* by recording a trail of edits in any SQL database. You connect it to an existing `*sql.DB`, pick a table name such as `audit_log`, and the package takes care of creating the table (if you want) and storing audit entries. Each entry captures the object that changed, who made the change, and the before/after values, so you can later review activity history or troubleshoot issues.

## Features
- **One-line setup**: Call `NewStore(...)` to bootstrap the audit table and get a ready-to-use store.
- **Human-readable audit entries**: `NewRecord()` gives you chainable setters for object type, object ID, author, and JSON snapshots of old/new values.
- **Flexible queries**: `NewRecordQuery()` lets you filter by user, object, or time range and supports pagination and sorting.
- **Driver-aware SQL**: Under the hood, Auditstore generates SQL that works with popular drivers like SQLite, PostgreSQL, MySQL, and SQL Server.
- **Optional verbose logging**: Flip on `EnableDebugMode(true)` to see the SQL statements being executed.

## Installation
```bash
go get github.com/dracory/auditstore
```
Ensure your project uses Go `1.26.3` or newer (specified in `go.mod`).

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
- **AuditGet(id string) (RecordInterface, error)**: Retrieves a single record or `nil, nil` if not found.
- **AuditList(query RecordQueryInterface) ([]RecordInterface, error)**: Returns all records matching the query.
- **AuditCount(query RecordQueryInterface) (int64, error)**: Counts rows for the given query.
- **AuditDelete(id string) error**: Removes the matching row.
- **AutoMigrate() error**: Creates the audit table if absent (deprecated alias for `MigrateUp`).

### NewStoreOptions (`store_implementation.go`)
- **DB**: Required `*sql.DB` connection. The driver name is auto-detected to select the correct SQL dialect.
- **AuditTableName**: Required table name (e.g., `audit_log`).
- **AutomigrateEnabled**: Auto-run `MigrateUp()` during construction.
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
`RecordQueryInterface` (`record_query_interface.go`) exposes fluent setters. The concrete implementation in `record_query_implementation.go` validates inputs and builds a `neat.Query` via `ToQuery(db)`.

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

Audit records are immutable — once written they are never updated. The table schema is declared in `MigrateUp` using the `neat` schema builder (`github.com/dracory/neat`), which generates dialect-appropriate DDL for SQLite, PostgreSQL, MySQL, and SQL Server.

## Debugging & Logging
- **Enable debug mode**: `store.EnableDebugMode(true)` switches the logger to debug level using `slog.NewTextHandler` with `LevelDebug`.
- **SQL output**: When debug is on, SQL statements and parameters are printed before execution across CRUD methods.

## Testing
`store_test.go` and `record_query_test.go` contain integration-style tests using `modernc.org/sqlite` in-memory databases:
- **`TestStoreAuditCreate`**: Ensures IDs are assigned and records persist.
- **`TestStoreAuditGet`**, **`TestStoreAuditGetNotFound`**: Confirm retrieval and correct nil return on missing records.
- **`TestStoreAuditList`**, **`TestStoreAuditCount`**, **`TestStoreAuditDelete`**: Cover the primary Store operations.
- **`TestRecordQueryToQuery`**: Exercises filtering, pagination, date ranges, and ordering end-to-end.

Run the suite with:
```bash
go test ./...
```

## Development Notes
- Ensure the target database driver is registered with `database/sql` (e.g., import `_ "modernc.org/sqlite"`).
- `AuditCreate` accepts any `RecordInterface` implementation — no type assertion is performed.
- Timestamps are stored as `time.Time` (UTC) so the ORM handles dialect-specific serialisation correctly.

## Dependencies
Auditstore relies on:
- **`github.com/dracory/neat`** for ORM operations and schema management.
- **`github.com/dracory/uid`** for unique ID generation.
- **`github.com/dromara/carbon/v2`** for time handling.
- **`modernc.org/sqlite`** in tests as an embedded driver; replace or supplement with your chosen driver in production.
