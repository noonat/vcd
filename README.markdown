# Vehicle Configuration Data Archive

Allows users to paste their [Captain Forever](http://captainforever.com) vessel exports to share with others.

## Development

vcd is written in [Go](https://golang.org) and uses [MySQL 8](https://www.mysql.com)
for the database. You can run it using docker-compose by running:

    docker-compose up

You can run unit tests with:

    go test ./...
