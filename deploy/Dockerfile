FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /portal ./cmd/portal

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /portal /portal
COPY specs /specs
USER nobody
ENV PORTAL_SPECS_DIR=/specs
ENV PORTAL_HTTP_ADDR=:8111
EXPOSE 8111
ENTRYPOINT ["/portal"]
