package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/marcoths/bookr/internal/data"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

var ErrFail = errors.New("FAIL")

func NewBookCmd(db *bbolt.DB) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "BOOK",
		Short: "Book seat(s) on a flight",
		Long:  `Book seat(s) on a flight by providing the seat number and number of consecutive seats to book as arguments.`,
		Args:  cobra.ExactArgs(2),
		RunE:  runBook(db),
	}
	return cmd
}

func runBook(db *bbolt.DB) RunEfn {
	return func(cmd *cobra.Command, args []string) error {
		seatArg := args[0]
		s, err := data.FromString(seatArg)
		if err != nil {
			return err
		}

		qty, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		available, err := getAvailableSeats(db, s, uint(qty))
		if err != nil {
			return ErrFail
		}

		if err := book(db, available); err != nil {
			return ErrFail
		}
		fmt.Println("SUCCESS")
		return nil
	}
}

func book(db *bbolt.DB, seats []*data.Seat) error {
	for _, s := range seats {
		if err := data.UpdateStatusTo(db, s.ID, data.Booked); err != nil {
			return err
		}
	}
	return nil
}

func getAvailableSeats(db *bbolt.DB, s *data.Seat, qty uint) ([]*data.Seat, error) {
	if err := data.IsValidQuantity(qty); err != nil {
		return nil, err
	}

	availableSeats, err := data.FindConsecutive(db, s.Row, s.Number, qty, data.Available)
	if err != nil {
		return nil, err
	}
	return availableSeats, nil
}
