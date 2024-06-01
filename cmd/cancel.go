package cmd

import (
	"fmt"
	"strconv"

	"github.com/marcoths/bookr/internal/data"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

func NewCancelCmd(db *bbolt.DB) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "CANCEL",
		Short: "Cancel seat(s) on a flight",
		Long:  `Cancel booked seat(s) on a flight by providing the seat number and number of consecutive seats to cancel as arguments.`,
		Args:  cobra.ExactArgs(2),
		RunE:  runCancel(db),
	}
	return cmd
}

func runCancel(db *bbolt.DB) RunEfn {
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

		booked, err := getBookedSeats(db, s, uint(qty))
		if err != nil {
			return err
		}
		if err := cancel(db, booked); err != nil {
			return err
		}

		fmt.Println("SUCCESS")
		return nil
	}
}

func getBookedSeats(db *bbolt.DB, s *data.Seat, qty uint) ([]*data.Seat, error) {
	if err := data.IsValidQuantity(qty); err != nil {
		return nil, err
	}

	booked, err := data.FindConsecutive(db, s.Row, s.Number, qty, data.Booked)
	if err != nil {
		return nil, err
	}
	return booked, nil
}

func cancel(db *bbolt.DB, seats []*data.Seat) error {
	for _, s := range seats {
		if err := data.UpdateStatusTo(db, s.ID, data.Available); err != nil {
			return err
		}
	}
	return nil
}
