APP=main

clear:
	rm -f ${APP} || true

build: clear
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./main ./url_shortener/cmd/url_shortener
