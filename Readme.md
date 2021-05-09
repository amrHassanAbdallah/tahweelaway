# Tahweelaway
> Centralized place to control your money

Moving money from your bank account to your friends.

## Getting started
1. install [golang](https://golang.org/), [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate), [sqlc](https://github.com/kyleconroy/sqlc#installation)
1. start a postgresql instance locally
1. migrate the db (make sure to replace the values in the connection string with yours)
   ```
   $ migrate -path persistence/migration -database "postgresql://root:secret@localhost:5432/tahweelaway?sslmode=disable" -verbose up
   ```
1. run the app
   ```
    $ make generate
    $ make build
    $ ./bin/app --postgresql-connection=postgresql://root:secret@localhost:5432/tahweelaway?sslmode=disable
   ```
### Or using docker-compose
1. make sure that you have docker, docker-compose installed
1. run
   ```
   $ docker-compose up
   ```

