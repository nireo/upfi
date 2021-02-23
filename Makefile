upfi: clean
	go get -u github.com/valyala/fasthttp
	go get -u gorm.io/gorm
	go get -u gorm.io/driver/postgres
	go get -u github.com/dgrijalva/jwt-go
	go get -u github.com/buaazp/fasthttprouter
	go get -u github.com/satori/go.uuid
	go build main.go -o upfi

clean:
	rm -rf upfi
