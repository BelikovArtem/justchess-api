FROM golang:1.22.5
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l" -o /bin/api/api main/main.go

# First port is listened by API, second port by Delve remove debugger.
EXPOSE 3502 4502 
COPY ./.env /bin/api/
COPY ./pkg/db/schema.sql /bin/api/
CMD ["/go/bin/dlv", "--listen=:4502", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "exec", "/bin/api/api"]