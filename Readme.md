# Tahweelaway
> Centralized place to control your money

Moving money from your bank account to your friends.

## Getting started
### Manually
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
### Using docker-compose
1. make sure that you have docker, docker-compose installed
1. run
   ```
   $ docker-compose up
   ```

### Check the API
Use this [file](https://github.com/amrHassanAbdallah/tahweelaway/blob/master/api/api.yml) content and paste it inside this [viewer](https://editor.swagger.io/)


## Features
* [ ] Authorize user
* [x] Create user
* [x] Get user
* [x] Create bank
* [ ] List banks
* [x] Transfer from bank to account
* [x] Transfer from account to account
* [ ] List transfers
* [ ] Refactor transfer domain to contain the currency 

Maybe will add more depending on this [design document](https://drive.google.com/file/d/185Y3opZoWqQNmEuNZBXyayxWA-3_i_kb/view?usp=sharing)
  
