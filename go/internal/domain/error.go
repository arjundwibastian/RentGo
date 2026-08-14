package domain

import "errors"

var ErrNotFound = errors.New("ERROR: data not found")
var ErrEmailRegistered = errors.New("ERROR: Email already registered!")
var InvalidEmailorPassword = errors.New("ERROR: Invalid email or password")
var ErrvehicleUnavailable = errors.New("ERROR: vehicle fully booked for selected date")
var ErrInvalidBookingUserID = errors.New("Error: user is not allowed to cancel this booking")
var ErrInvalidConfirmInput = errors.New("Error:Invalid cancel confirmation input")
var ErrNotEnoughBalance = errors.New("Error: your balance is not enough to make this booking! please top up to increase your balance.")
var ErrInvalidBookingCancel = errors.New("Error: your booking cant be cancelled because booking is completed or already cancelled")