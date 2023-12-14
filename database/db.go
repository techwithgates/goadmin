package database

import "github.com/jackc/pgx/v5/pgxpool"

var Db *pgxpool.Pool

// statement to get table names
var RetrieveTableStmt = `
	SELECT table_name
	FROM information_schema.tables
	WHERE table_schema='public'
	ORDER BY table_name
`

// statement to get the tables whose table schema is public
var TableVerifyStmt = `
	SELECT table_name
	FROM information_schema.tables
	WHERE table_name=$1 AND table_schema='public'
`

// statement to get the PK field of a table by default constraint name
var FindPkStmt = `
	SELECT column_name, constraint_name
	FROM information_schema.key_column_usage
	WHERE table_name=$1 AND constraint_name=$2
`

// get the PK field type (SERIAL or other else)
var PkTypeStmt = `
	SELECT column_name, column_default
	FROM information_schema.columns
	WHERE table_name=$1 AND column_name=$2
`

// statement to construct HTML form field definitions
var GetMetaDataStmt = `
	SELECT column_name, data_type, is_nullable
	FROM information_schema.columns
	WHERE table_name=$1 AND column_name!=$2
	ORDER BY ordinal_position
`

// statement to a list of PKs
var RetrievePksStmt = `SELECT %s FROM %s ORDER BY %s LIMIT %d OFFSET %d`

// statement to insert dynamic data to dynamic column fields
var InsertDataStmt = "INSERT INTO %s (%s) VALUES (%s)"

// statement to update a single data record
var UpdateDataStmt = "UPDATE %s SET %s WHERE %s=$1"

// statement to get a single value of a column field
var GetFieldStmt = `SELECT %s FROM %s WHERE %s=$1`

var GetTotalObjStmt = `SELECT COUNT(*) FROM %s`

// statement to delete a single data record
var DeleteDataStmt = `DELETE FROM %s WHERE %s IN (%s)`

// statement to distinguish char and file field
var GetColNumber = `
	SELECT attnum
	FROM pg_attribute
	WHERE attrelid='%s'::regclass AND attname='%s'
`

// statement to distinguish char and file field
var GetCommentStmt = `SELECT col_description('%s'::regclass, %v)`
