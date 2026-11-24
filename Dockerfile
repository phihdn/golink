FROM golang:1.27.4-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/golink ./cmd/api

FROM alpine:3.22
COPY --from=build /out/golink /usr/local/bin/golink
EXPOSE 8080
ENTRYPOINT ["golink"]
