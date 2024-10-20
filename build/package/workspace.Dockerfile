FROM golang:1.23-alpine AS build

WORKDIR /todennus-backend

RUN apk add -U --no-cache ca-certificates

COPY ./todennus-backend/go.mod .
COPY ./todennus-backend/go.sum .

RUN go mod download

COPY . /

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /todennus ./cmd/main.go

FROM scratch

WORKDIR /

COPY --from=build /todennus /

EXPOSE 8080

ENTRYPOINT [ "/todennus", "rest", "--env", ""]
