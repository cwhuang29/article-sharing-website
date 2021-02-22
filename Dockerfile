FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y wget build-essential
 
RUN ["wget", "https://dl.google.com/go/go1.15.6.linux-amd64.tar.gz"]
RUN ["tar", "-C", "/usr/local", "-xzf", "go1.15.6.linux-amd64.tar.gz"]
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app
COPY . .
RUN go build -o web /app/cmd/main.go

ENTRYPOINT ["/app/web"]
