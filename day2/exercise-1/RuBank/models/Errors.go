package models

import "errors"

var ErrUsersNotFound = errors.New("users not found")
var ErrInvalidTransfer = errors.New("invalid transfer")
var ErrSmallAmount = errors.New("the amount is too small")
var ErrYourself = errors.New("translation to yourself")
var ErrSuspicious = errors.New("suspicious translation")
var ErrPrintID = errors.New("Error Print ID")
