FROM golang:1.13-alpine3.11 as builder

RUN apk --no-cache add gcc g++ make ca-certificates git
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main ./src/.

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# STAGE 2
# generate clean, final image for end users
#
FROM alpine:3.11.3

ENV TZ=Asia/Jakarta
RUN apk --update add tzdata ca-certificates && \
    update-ca-certificates 2>/dev/null || true \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# copy golang binary into container
COPY --from=builder /dist/ .

EXPOSE 8080

# executable
ENTRYPOINT [ "./main" ]
