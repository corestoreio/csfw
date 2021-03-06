package pseudo

import (
	"sync/atomic"
	"time"

	"github.com/oklog/ulid"
	"github.com/satori/go.uuid"
)

var lowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")
var upperLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numeric = []rune("0123456789")
var specialChars = []rune(`!'@#$%^&*()_+-=[]{};:",./?`)
var hexDigits = []rune("0123456789abcdef")

func (s *Service) text(atLeast, atMost int, allowLower, allowUpper, allowNumeric, allowSpecial bool) string {
	allowedChars := []rune{}
	if allowLower {
		allowedChars = append(allowedChars, lowerLetters...)
	}
	if allowUpper {
		allowedChars = append(allowedChars, upperLetters...)
	}
	if allowNumeric {
		allowedChars = append(allowedChars, numeric...)
	}
	if allowSpecial {
		allowedChars = append(allowedChars, specialChars...)
	}

	result := []rune{}
	nTimes := s.r.Intn(atMost-atLeast+1) + atLeast
	for i := 0; i < nTimes; i++ {
		result = append(result, allowedChars[s.r.Intn(len(allowedChars))])
	}
	return string(result)
}

// Password generates password with the length from atLeast to atMOst charachers,
// allow* parameters specify whether corresponding symbols can be used
func (s *Service) Password(atLeast, atMost int, allowUpper, allowNumeric, allowSpecial bool) string {
	return s.text(atLeast, atMost, true, allowUpper, allowNumeric, allowSpecial)
}

// SimplePassword is a wrapper around Password,
// it generates password with the length from 6 to 12 symbols, with upper characters and numeric symbols allowed
func (s *Service) SimplePassword() string {
	return s.Password(6, 12, true, true, false)
}

// Color generates color name
func (s *Service) Color() string {
	return s.lookup(s.o.Lang, "colors", true)
}

// DigitsN returns n digits as a string
func (s *Service) DigitsN(n int) string {
	digits := make([]rune, n)
	for i := 0; i < n; i++ {
		digits[i] = numeric[s.r.Intn(len(numeric))]
	}
	return string(digits)
}

// Digits returns from 1 to 5 digits as a string
func (s *Service) Digits() string {
	return s.DigitsN(s.r.Intn(5) + 1)
}

func (s *Service) hexDigitsStr(n int) string {
	var num []rune
	for i := 0; i < n; i++ {
		num = append(num, hexDigits[s.r.Intn(len(hexDigits))])
	}
	return string(num)
}

// HexColor generates hex color name
func (s *Service) HexColor() string {
	return s.hexDigitsStr(6)
}

// HexColorShort generates short hex color name
func (s *Service) HexColorShort() string {
	return s.hexDigitsStr(3)
}

// ID returns a never seen before unique ID.
func (s *Service) ID() uint64 {
	return atomic.AddUint64(s.id, 1)
}

// UUID returns a UUID v4.
func (s *Service) UUID() []byte {
	u := uuid.Must(uuid.NewV4())
	return u.Bytes()
}

// UUID returns a UUID v4 as string.
func (s *Service) UUIDString() string {
	u := uuid.Must(uuid.NewV4())
	return u.String()
}

// ULID returns an ULID which is a 16 byte Universally Unique Lexicographically
// Sortable Identifier.
func (s *Service) ULID() ulid.ULID {
	return ulid.MustNew(ulid.Timestamp(time.Now()), s.ulidEntropy)
}

// Intn returns a non-negative pseudo-random int.
func (s *Service) Intn(n int) int {
	return s.r.Intn(n)
}

func (s *Service) Float64() float64 {
	return float64(s.r.Intn(10000)) + float64(s.r.Uint64())/1e12
}
