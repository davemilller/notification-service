# build stage
FROM golang:1.17-buster as builder

WORKDIR /app

COPY app/go.* /app
COPY app/vendor* /app/vendor

# creates a script to run go mod download only if vendor directory does not exist
RUN echo "#!/bin/sh\n" \
    "if [ -d '/app/vendor' ] ; then\n" \
    "   echo '--> using vendored dependencies' ;\n" \
    "else\n" \
    "   echo '--> downloading dependencies...\n--> speed up the build by running go mod vendor!' \n" \
    "   go mod download ;\n" \
    "fi\n" > script.sh

RUN chmod +x script.sh
RUN ./script.sh

COPY app /app

# build binary
RUN go build -o /notification-service ./cmd/main.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates libc6-compat

COPY --from=builder /notification-service .

CMD ["/notification-service"]
