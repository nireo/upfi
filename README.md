# upfi: file management server

[![Go Report Card](https://goreportcard.com/badge/github.com/nireo/upfi)](https://goreportcard.com/report/github.com/nireo/upfi)

## Goal

The goal of the project is to create an easy to setup file hosting service. The idea is that anyone with a linux computer can setup a upfi-instance!


## Setup

You will need to create a database and configure an environment variables file. Here is a example of the `.env` file. All of the fields below must be added to the `.env` for the service to work.

```go
#.env
db_name=upfi
db_port=5432
db_host=localhost
port=8080
db_user=postgres
root_dir=/home/username/go/src/github.com/nireo/upfi/
```

The root dir is there since I found some problems with relative file paths. Such that the project uses a util function which appends the 'root_dir' variable to all of the paths.

To use a different database than postgres, check out the [documentation](https://gorm.io/docs/connecting_to_the_database.html) of gorm.


Now you can just run the app.

```
go run main.go
```

## Contributions

Anyone can contribute to the project by creating a pull request!
