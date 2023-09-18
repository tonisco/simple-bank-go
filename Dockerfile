# build stage
FROM golang:1.21.1-alpine3.17 AS Builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

RUN apk add curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz



# RUN stage 
FROM alpine:3.17

WORKDIR /app

COPY --from=Builder /app/main .

COPY --from=Builder /app/migrate ./migrate

COPY app.env .

COPY start.sh .

COPY wait-for.sh .

COPY db/migration ./migration

EXPOSE 8080

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
