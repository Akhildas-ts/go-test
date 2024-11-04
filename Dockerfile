# syntax=docker/dockerfile:1
FROM golang:1.23.2 AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /test 

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /test /test

COPY .env /app/.env

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/test"]