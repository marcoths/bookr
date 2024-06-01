package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SeatStatus int

const (
	Available SeatStatus = iota
	Booked
)

func (s SeatStatus) String() string {
	return [...]string{"available", "booked"}[s]
}

var OutOfBoundsError = errors.New("seat out of bounds")

type Seat struct {
	ID     string     `json:"id"`
	Row    string     `json:"row"`
	Number uint       `json:"number"`
	Status SeatStatus `json:"status"`
}

type Seats []*Seat

// FromString validates a seat string and returns a Seat object
func FromString(s string) (*Seat, error) {
	if !isInBounds(s) {
		return nil, fmt.Errorf("%s: %w", s, OutOfBoundsError)
	}
	parts := strings.Split(s, "")
	row := strings.ToUpper(parts[0])
	number, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: %w", err, OutOfBoundsError)
	}

	if err := IsValidNumber(uint(number)); err != nil {
		return nil, fmt.Errorf("seat number %d: %w", number, err)
	}

	if !isValidRow(row) {
		return nil, fmt.Errorf("%s: %w", row, OutOfBoundsError)
	}
	return &Seat{ID: s, Row: row, Number: uint(number)}, nil
}

func IsValidNumber(number uint) error {
	if !(number > 0 && number <= 7) {
		return OutOfBoundsError
	}
	return nil
}

func isInBounds(str string) bool {
	return len(str) > 1 && len(str) < 3
}

func isValidRow(row string) bool {
	return row >= "A" && row <= "T"
}
