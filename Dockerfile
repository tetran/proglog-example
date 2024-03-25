FROM golang:1.21-alpine AS build
WORKDIR /go/src/proglog-example
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/proglog-example ./cmd/proglog-example

# q: what is scratch?
# a: scratch is a special Docker image that is completely empty. It is often used as a base image for building minimal Docker images.
FROM scratch
COPY --from=build /go/bin/proglog-example /bin/proglog-example
ENTRYPOINT [ "/bin/proglog-example" ]
