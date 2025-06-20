FROM golang:alpine

WORKDIR /practicum
COPY . .
RUN go mod download

RUN go build -o main .

EXPOSE 8080
CMD ["./main"]

