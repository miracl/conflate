FROM golang:latest
WORKDIR /src
COPY . .
RUN go build -o /usr/local/bin/conflate .
ENTRYPOINT ["conflate"]
