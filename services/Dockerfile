FROM golang:1.16

WORKDIR /app

COPY ./ /app/

RUN go build -o kim

EXPOSE 8000

ENTRYPOINT $entrypoint