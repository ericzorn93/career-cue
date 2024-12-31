FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.work .
COPY go.work.sum .
COPY ./libs ./libs
COPY ./apps/services ./apps/services
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -o /bin/accounts-graphql ./apps/services/accounts-graphql/cmd/server/main.go

FROM alpine:3.20
RUN apk update --no-cache bash curl
WORKDIR /app
COPY --from=builder /bin/accounts-graphql /bin/accounts-graphql
EXPOSE 3000
CMD [ "/bin/accounts-graphql" ]
