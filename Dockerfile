FROM golang:latest AS builder
ADD .. /url_shortener
WORKDIR /url_shortener
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/url-shortener/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY /public ./public
COPY --from=builder url_shortener/main ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]
EXPOSE 8080