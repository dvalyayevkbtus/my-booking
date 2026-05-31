# my-booking

This is a sample Golang HTTP service.

## Endpoints

**App is running on port 8080**

There are 3 endpoints:

- GET `/client` - To get all clients
- POST `/client` - To create new client
    Body:
    ```json
    {
        "fullName": "John Snow"
    }
    ```
- GET `/booking` - To get all bookings
- POST `/booking` - To register new booking
    Body:
    ```json
    {
        "hotelName": "Test Plaza",
        "price": "10000.00",
        "currency": "KZT",
        "clientId": 1
    }
    ```
- GET `/booking/{id}` - To get booking by id

## How to build

### Host

1. Check if you have go `go version` or install it
2. Run `go build .`. The file will be located in this dir `my-booking`

### Dockerfile

1. Use `golang:1.24.2` as a base image for build
2. Copy `go.mod` and `go.sum` and run `go mod download`
3. Copy all files inside and run `go build .` or `go build -o app`
4. Prepare a running image. You can use `FROM alpine`.
5. Copy binary from build to run and configure entrypoint to just run your binary

## How to run

### Prepare a config

This app needs a PostgreSQL running and payment service (https://github.com/dvalyayevkbtu/payment).
Create a file `my-booking.json`:
    
```json
{
    "db": {
        "host": "<host of postgres>",
        "port": "<port of postgres>",
        "name": "<database name>",
        "user": "<database username>",
        "password": "<database password>"
    },
    "payment": {
        "url": "http://<host>:<port>"
    }
}
```

### Configure environment

You need to setup `BOOKING_CONFIG_PATH` environment variable with a full path to your config file.

### Running

You need just run

## How to test

You can execute following commands:
```shell
curl -X POST -d '{"fullName": "Test"}' http://<host>:<port>/client
curl -X POST -d '{"hotelName": "Test", "price": "10000.00", "currency": "USD", "clientId": 1}' http://<host>:<port>/booking
```

After some time, you need to execute:
```shell
curl http://<host>:<port>/booking
```

And you need to see: `"paid":true`
