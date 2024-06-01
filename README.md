# BOOKR

CLI tool to manage airplane bookings.

## Installation

```bash
$ make build
```
This will create binaries for Linux and MacOS in the `bin` directory.

## Usage

```bash
bookr [command]

Available Commands:
  book        Book seat(s) on a flight
  cancel      Cancel seat(s) on a flight
  help        Help about any command
```

### Booking a seat

Book seat(s) on a flight by providing the seat number and number of consecutive seats to book as arguments. The airplane
has 20 rows and 8 seats per row. The seat number is a combination of the row number and the seat letter. For example, the
first seat in the first row is `A0` and the last seat in the last row is `T7`.
```bash
bookr book A1 2 # Book 2 consecutive seats starting from A1
```

### Cancelling a seat

Cancel seat(s) on a flight by providing the seat number and number of consecutive seats to cancel as arguments.
```bash
bookr cancel A1 2 # Cancel 2 consecutive seats starting from A1
```

### Running tests

```bash
$ make test
```
