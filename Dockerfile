FROM golang:1.22.2-alpine

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o mintyplex-api .

EXPOSE 8081

CMD ["./mintyplex-api"]
