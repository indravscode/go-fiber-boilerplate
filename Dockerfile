FROM golang:1.25 AS build

WORKDIR /app
COPY . .
RUN go clean --modcache
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build src/main.go

FROM alpine:latest

RUN apk add --no-cache curl tzdata

WORKDIR /root
COPY --from=build /app/main .
COPY --from=build /app/.env .

EXPOSE 3000
CMD ["./main"]
