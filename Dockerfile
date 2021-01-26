FROM golang:1.15.2-alpine3.12 as go-builder
ARG TARGETPLATFORM
WORKDIR /work
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN sh ./scripts/build.sh sensoapid $TARGETPLATFORM

FROM alpine:3.12.0
ARG TARGETPLATFORM
COPY --from=go-builder /work/bin/$TARGETPLATFORM /app/bin
EXPOSE 8080
CMD ["/app/bin/sensoapid"]