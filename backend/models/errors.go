package models

import "fmt"

type CError error

var (
	ErrCartNotYours  = fmt.Errorf("cart is not yours")
	ErrCartInProcess = fmt.Errorf("cart is in process")
	ErrCartNotFound  = fmt.Errorf("cart not found")
)
