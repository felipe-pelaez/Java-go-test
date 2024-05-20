FROM golang:alpine3.19 AS golang

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server .

FROM gcr.io/distroless/static-debian12

COPY --from=golang /server .

CMD ["/server"]
