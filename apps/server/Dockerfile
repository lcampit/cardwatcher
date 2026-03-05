FROM golang:1.25.7-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY gen/ ./gen/
COPY apps/server/cmd/ ./apps/server/cmd/
COPY apps/server/internal/ ./apps/server/internal/

RUN go build -o main apps/server/cmd/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
RUN apk update && apk --no-cache add tzdata ca-certificates
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


