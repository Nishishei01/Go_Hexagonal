FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

FROM alpine:3.18

COPY --from=build /app/server /server

EXPOSE 3000

CMD ["/server"]