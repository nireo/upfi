# upfi: file management server

> upfi is a minimal version of google drive.

## Setup

Upfi is built with golang, so you will need to install it. Also upfi has a few go dependencies:

```
# handling http sessions via cookies
go get -u github.com/gorilla/sessions

# ORM library to work easily with postgresql
go get -u github.com/jinzhu/gorm
```

You will need to create a database:

```go
// main.go
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=upfi sslmode=disable")
	if err != nil {
		panic(err)
	}
	models.MigrateModels(db)
	defer db.Close()
	lib.SetDatabase(db)
```

To use a different database than postgres, check out the [documentation](https://gorm.io/docs/connecting_to_the_database.html) of gorm.

Then just type:

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
