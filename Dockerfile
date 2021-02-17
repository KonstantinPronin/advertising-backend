FROM golang:1.15.3-alpine
COPY . /home/app
WORKDIR /home/app
RUN go build -o build/app ./cmd/main/main.go
CMD ./build/app