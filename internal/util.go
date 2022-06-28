package internal

import "time"

// Now returns Todays Date with truncated time parts
// mimicks time.Now() from "time" package, for testing
func Now() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// ToStringPtr returns pointer to the given string
func ToStringPtr(s string) *string {
	return &s
}
