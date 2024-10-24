FROM golang:1.23-alpine AS build

WORKDIR /todennus-backend

RUN apk add -U --no-cache ca-certificates

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /todennus ./cmd/main.go

FROM scratch

WORKDIR /

COPY --from=build /todennus /
COPY --from=build /todennus-backend/template /template
COPY --from=build /todennus-backend/docs /docs

EXPOSE 8080 8081 8083

ENTRYPOINT [ "/todennus", "--env", ""]
