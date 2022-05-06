FROM golang:latest
WORKDIR /usr/src/app

COPY . .

RUN go get
RUN go install github.com/pilu/fresh@latest

EXPOSE 8000