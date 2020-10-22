package types

import (
	"regexp"
	"strings"
)

var cellphoneRegex = regexp.MustCompile(`^09\d{9}$`)

// Cellphone is the data-type for cellphone numbers.
type Cellphone string

// IsValid returns true if the cellphone number is valid.
func (p Cellphone) IsValid() bool {
	return cellphoneRegex.MatchString(p.Standard().String())
}

// Standard converts the cellphone number to standard format.
func (p Cellphone) Standard() Cellphone {
	switch {
	case strings.Index(p.String(), "0098") == 0:
		return "0" + p[4:]
	case strings.Index(p.String(), "+98") == 0:
		return "0" + p[3:]
	case strings.Index(p.String(), "98") == 0:
		return "0" + p[2:]
	case strings.Index(p.String(), "9") == 0:
		return "0" + p
	default:
		return p
	}
}

// Masked returns a masked cellphone number.
func (p Cellphone) Masked() string {
	std := p.Standard().String()
	// 0912xxxxx67
	return std[:4] + "xxxxx" + std[9:]
}

func (p Cellphone) String() string {
	return string(p)
}
