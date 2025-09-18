package database

import "errors"

var ErrDatabaseIsDirty = errors.New("database is dirty, please resolve the issue and try again")

var ErrCannotRunMigrations = errors.New("cannot run migrations")
