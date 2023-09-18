FROM golang:1.19 AS builder
WORKDIR /l0wb
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./app ./main.go

FROM alpine:latest
WORKDIR /noteservice
COPY ./.env .
COPY ./template ./template
COPY --from=builder /l0wb/app .
ENTRYPOINT [ "/noteservice/app" ]