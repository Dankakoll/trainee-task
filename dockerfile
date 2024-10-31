# docker build -t eps_test_task .
FROM golang:1.23.2
COPY . .
RUN go test -v
