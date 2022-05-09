# syntax=docker/dockerfile:1

FROM golang:1.18.1-alpine3.14

WORKDIR /app

COPY . .
RUN go mod download
RUN make swag
RUN make build

EXPOSE 8080

RUN docker-compose up
CMD [ "main" ]


