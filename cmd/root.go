package cmd

import (
	"fmt"

	"github.com/marcoths/bookr/internal/bucket"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

type RunEfn func(cmd *cobra.Command, args []string) error

func NewRootCmd(db *bolt.DB) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "bookr",
		Short: "A simple tool for managing airline bookings",
		Long: `bookr is a simple tool for managing airline bookings. 
It allows you to book seats on a flight and cancel bookings.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(NewBookCmd(db))
	rootCmd.AddCommand(NewCancelCmd(db))

	return rootCmd
}

func Register(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Create all the buckets except auth, it will be created in setParameters()
		name := bucket.Seats.Name()
		if _, err := tx.CreateBucketIfNotExists(name); err != nil {
			return fmt.Errorf("create bucket %s: %w", name, err)
		}
		return nil
	})
}
