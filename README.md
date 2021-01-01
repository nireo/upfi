# upfi: file management server

> upfi is a minimal version of google drive.

## Goal

The goal of the project is to create an easy to setup file hosting service. The idea is that anyone with a linux computer can setup a upfi-instance!

## Setup

Upfi is built with golang, so you will need to install it. Also upfi has a few go dependencies. These all are downloaded with the `Makefile` when you run make.

Example of the Makefile:

```
upfi: clean
	go get -u github.com/valyala/fasthttp
	go get -u gorm.io/gorm
	go get -u gorm.io/driver/postgres
	go get -u github.com/gorilla
	go get -u github.com/dgrijalva/jwt-go
	go get -u github.com/gorilla/sessions
	go get -u github.com/buaazp/fasthttprouter
	go get -u github.com/satori/go.uuid
	go build

clean:
	rm -rf upfi
```

You will need to create a database and configure an environment variables file. Here is a example of the `.env` file. All of the fields below must be added to the `.env` for the service to work.

```go
#.env
db_name=upfi
db_port=5432
db_host=localhost
port=8080
db_user=postgres
```

To use a different database than postgres, check out the [documentation](https://gorm.io/docs/connecting_to_the_database.html) of gorm.

To build the program just type:

```
make
```

Or just to the run the service:

```
go run main.go
```

## Reset information

There is the functionality to remove all files and all user data. This can be done by simply running:

```
go run main.go reset_information
```

After executing this command the program will just end. **NOTE this cannot be reversed.**

## Contributions

Anyone can contribute to the project by creating a pull request!
