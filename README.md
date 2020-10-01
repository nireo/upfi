# upfi: file management server

> upfi is a minimal version of google drive.

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

## Optimized API (WIP)

There is currently an experimental handler rewrite going on which uses the [fasthttp](https://github.com/valyala/fasthttp) package instead of the default `net/http`. Some features in the handlers are experimental and they are in the middle of being optimized and being more secure. Since it's experimental some of the features might not work!

### Running the optimized api

```
# Without the api flag the default api value will be 'default', which uses the net/http package
go run main.go -api=optimized
```

In the future encrypting files is planned for the optimized api. Also in the future the optimized api might become the new default, but that needs some more integration!

## Contributions

Anyone can contribute to the project by creating a pull request!
