package db

import "fmt"

type ValidationError struct {
	Field  string
	Detail string
}

func (r *ValidationError) Error() string {
	return fmt.Sprintf("invalid_argument:%s:%s", r.Field, r.Detail)
}
