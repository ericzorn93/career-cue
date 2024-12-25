FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.work .
COPY go.work.sum .
COPY ./libs ./libs
COPY ./apps/services ./apps/services
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -o /bin/accounts-api ./apps/services/accounts-api/cmd/server/main.go

FROM alpine:3.20
RUN apk update --no-cache bash curl
WORKDIR /app
COPY --from=builder /bin/accounts-api /bin/accounts-api
EXPOSE 3000
EXPOSE 5000
CMD [ "/bin/accounts-api" ]
