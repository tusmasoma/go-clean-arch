package email

import (
	"errors"
	"strings"
)

// GetAddressPart extracts the part before '@' in the given email string
// and returns it as a string.
// If '@' is not present, it returns an error.
func GetAddressPart(email string) (string, error) {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid email address: missing '@' or multiple '@'")
	}
	return parts[0], nil
}
