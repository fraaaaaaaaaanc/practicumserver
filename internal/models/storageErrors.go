package models

import "errors"

// Ошибка сообщающая о конфликте данных в хранилище
var ErrConflictData = errors.New("data conflict, the resulting url already exists in the storage")

// Ошибка сообщающая о том, что в таблице нет таких данных
var ErrNoRows = errors.New("there is no such data in the table")

// Ошибка сообщающая о том, что данные были удалены
var ErrDeletedData = errors.New("this data has been deleted")

// Ошибка сообщающая о том, что тип данных передаваемых в контексте не верен
var ErrUserIDType = errors.New("context variable type mismatch")
