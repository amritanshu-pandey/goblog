FROM docker.io/library/golang:1.23

WORKDIR /src
COPY . /src/
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN go mod download
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o / ./...
RUN ls -lrtah /

FROM scratch
COPY --from=0 /goblog /bin/goblog
CMD ["/bin/goblog"]
