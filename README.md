# fetcher
Purpose of the application is to fetch responses from added urls with given interval and save them as a history.

Project comes with configuration file: _fetcher.yml_, which contains a database and api configuration.

To run this repository you can:
```
make binary
./fetcher
```
or
```
go run cmd/fetcher/main.go
```

Project comes with a Makefile that contains following recipes:

- `make binary` makes binary which can be run with `./fetcher`
- `make test` run tests
- `make cover` run tests with coverage results saved to file
- `make cover-total` same as above, but also print total coverage
- `make cover-html` runs test and open results in browser
