package data_test

import (
	"errors"
	"testing"

	"github.com/marcoths/bookr/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestSeats(t *testing.T) {

	tests := []struct {
		input    string
		expected *data.Seat
		wantErr  error
	}{
		{"A1", &data.Seat{ID: "A1", Row: "A", Number: 1}, nil},
		{"B2", &data.Seat{ID: "B2", Row: "B", Number: 2}, nil},
		{"C3", &data.Seat{ID: "C3", Row: "C", Number: 3}, nil},
		{"T7", &data.Seat{ID: "T7", Row: "T", Number: 7}, nil},
		{"Z9", nil, data.OutOfBoundsError},
		{"A0", &data.Seat{ID: "A0", Row: "A", Number: 0}, nil},
		{"A9", nil, data.OutOfBoundsError},
		{"", nil, data.OutOfBoundsError},    // empty string
		{"A", nil, data.OutOfBoundsError},   // only row
		{"1", nil, data.OutOfBoundsError},   // only number
		{"AA", nil, data.OutOfBoundsError},  // invalid number
		{"A10", nil, data.OutOfBoundsError}, // more than one digit for number
		{"A1B", nil, data.OutOfBoundsError}, // extra character
		{"123", nil, data.OutOfBoundsError}, // invalid row
		{"AB", nil, data.OutOfBoundsError},  // non-numeric number
		{"@1", nil, data.OutOfBoundsError},  // special character row
		{"A-1", nil, data.OutOfBoundsError}, // negative number
		{"AA", nil, data.OutOfBoundsError},  // double characters row
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := data.FromString(tt.input)
			assert.Equal(t, tt.expected, got, "FromString(%s) = %v, want %v", tt.input, got, tt.expected)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FromString(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}

}
