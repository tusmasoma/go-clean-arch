FROM golang:1.21.3

RUN apt-get update && apt-get install -y default-mysql-client && apt-get install -y postgresql-client

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

WORKDIR /app
RUN go mod download

RUN go install github.com/cosmtrek/air@v1.29.0

WORKDIR /app

COPY . .

COPY entrypoint.sh /usr/local/bin/

RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]

WORKDIR /app
CMD ["air", "-c", ".air.toml"]