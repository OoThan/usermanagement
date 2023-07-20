FROM golang:1.17
WORKDIR /go/src/app
COPY . .
RUN go build -o app
EXPOSE 7001
CMD ["./app"]