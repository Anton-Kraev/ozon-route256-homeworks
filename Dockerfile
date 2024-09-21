## Build

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.22 AS build

WORKDIR /cmd

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# copy all files
COPY . ./
# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app -v ./cmd

## Deploy
FROM alpine:3.20 AS final

WORKDIR /

COPY --from=build /bin/app /app
COPY ./configs /configs
COPY .env .env

EXPOSE 8080
