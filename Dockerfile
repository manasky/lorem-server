FROM golang:1.16 AS build

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -ldflags "-w -s" -o ./bin/lorem main.go

FROM alpine:3.15

WORKDIR /app

COPY --from=build /app/bin/lorem /app/lorem

USER 1000:1000

EXPOSE 8080
ENTRYPOINT ["./lorem"]