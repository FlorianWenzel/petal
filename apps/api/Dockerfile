FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/petal-api .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 1000 app
USER app
WORKDIR /app
COPY --from=build /out/petal-api /app/petal-api
EXPOSE 4000
ENTRYPOINT ["/app/petal-api"]
