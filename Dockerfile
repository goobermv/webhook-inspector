FROM golang:1.26.4-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/server ./cmd/server

FROM alpine:3.20
RUN adduser -D app
USER app
COPY --from=build /bin/server /server
COPY --from=build /app/internal/database/migrations /migrations
ENV MIGRATIONS_DIR=file:///migrations
EXPOSE 8080
ENTRYPOINT ["/server"]