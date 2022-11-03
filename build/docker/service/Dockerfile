FROM golang:1.19-alpine AS builder
WORKDIR /balancer-service
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY ./ ./
WORKDIR /balancer-service/cmd/app
RUN CGO_ENABLED=0 go build -o /bin/app


FROM alpine:latest 
RUN apk --update add ca-certificates

COPY --from=builder /bin/app /bin/app
COPY --from=builder /balancer-service/config/config.yaml /config/config.yaml
ENV PATH=/bin
EXPOSE 9096
CMD [ "/bin/app" ]