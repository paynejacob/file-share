FROM golang:1.16-alpine3.13 as gobuild
WORKDIR /file-share
COPY web web
COPY go.* ./
COPY main.go main.go
RUN GOOS=linux go build -o file-share main.go

FROM alpine:3.13
COPY --from=gobuild /file-share/file-share /usr/local/bin/file-share
EXPOSE 80
ENTRYPOINT ["/usr/local/bin/file-share"]