# Initial stage: download modules
FROM golang:1.13 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.13 as builder

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 myapp

RUN mkdir -p /paris
ADD . /paris
WORKDIR /paris

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 make build

# Final stage: Run the binary
FROM scratch

ENV PORT 8080
ENV DIAG_PORT 8181

# certificates to interact with other services
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# don't forget /etc/passwd from previous stage
COPY --from=builder /etc/passwd /etc/passwd
USER myapp

# and finally the binary
COPY --from=builder /paris/bin/paris /paris
EXPOSE $PORT
EXPOSE $DIAG_PORT

CMD ["/paris"]
