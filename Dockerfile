FROM golang:1.22

WORKDIR /app

COPY . /app

RUN go build -o task.exe ./cmd

CMD ["./task.exe", "input.txt"]