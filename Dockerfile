# build stage
FROM golang:1.21.1-alpine3.17 AS Builder

WORKDIR /app

COPY . .

RUN go build -o main main.go

# RUN stage 
FROM alpine:3.17

WORKDIR /app

COPY --from=Builder /app/main .
COPY app.env .

EXPOSE 8080

CMD [ "/app/main" ]
