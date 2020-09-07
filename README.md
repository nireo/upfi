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

## Contributions

Anyone can contribute to the project by creating a pull request!
