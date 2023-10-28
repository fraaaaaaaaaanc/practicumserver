package models

import "errors"

// ErrConflictData is an error indicating a data conflict, where the resulting URL already exists in the storage.
var ErrConflictData = errors.New("data conflict, the resulting url already exists in the storage")

// ErrNoRows is an error indicating that there is no such data in the table.
var ErrNoRows = errors.New("there is no such data in the table")

// ErrDeletedData is an error indicating that this data has been deleted.
var ErrDeletedData = errors.New("this data has been deleted")

// ErrUserIDType is an error indicating that the type of data passed in the context is incorrect.
var ErrUserIDType = errors.New("context variable type mismatch")
