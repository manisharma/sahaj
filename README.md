# Sahaj

Sahaj is an application designed to model a Parking Lot intended for parking vehicles.


## Configuration

It needs a `sample.json` file to load the initial configuration a Parking Lot should have like category of vehicles, their number of spots an parking Tarrif.


## Architecture

It uses `Factory Design Pattern` to create a Parking Lot based on `Type`.

In order for a `Parking Lot` to be considered a `Parking Lot` it has to satisfy a `Contract`

All implementations of the `Contract` & sensitive `Business Logics` are kept under `internal`

## Usage

use `make` instructions to `execute`, `test` the app.
We can test the code by using executing `main.go` or by changing the test cases.

```bash
# run all tests
make test 

# check test coverage
make coverage

# run the app
make run
```
