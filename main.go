package main

import (
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/marcoths/bookr/cmd"
	"github.com/marcoths/bookr/internal/data"
	bolt "go.etcd.io/bbolt"
)

var (
	//go:embed files/seed.json
	seed embed.FS
)

func main() {
	db, err := bolt.Open("bookr.db", 0600, &bolt.Options{Timeout: 200 * time.Millisecond})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	jsonData, err := seed.ReadFile("files/seed.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := data.Register(db); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := data.Seed(db, jsonData); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.NewRootCmd(db).Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

}
