FROM golang:1.18-alpine

RUN apk add --no-cache git

WORKDIR /app/cfp-api

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./build/cfp-api .

CMD [ "./build/cfp-api" ]