FROM golang:1.23-alpine AS builder
RUN apk update --no-cache bash curl
WORKDIR /app
COPY ../../. .
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -o /bin/accounts-worker ./apps/services/accounts-worker/cmd/server/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /bin/accounts-worker .
CMD [ "./accounts-worker" ]