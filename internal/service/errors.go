package service

import "errors"

var ErrOrderAlreadyExistsForAnotherUser = errors.New("order already exists for another user")
var ErrOrderAlreadyExists = errors.New("order already exists")
var ErrNotFound = errors.New("not found")
