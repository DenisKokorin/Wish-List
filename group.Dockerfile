FROM golang:1.24.2 AS group-builder
WORKDIR /app
COPY . .
RUN cd cmd/group && CGO_ENABLED=0 GOOS=linux go build -o group-service main.go

FROM alpine AS group
WORKDIR /
COPY --from=group-builder /app/cmd/group/group-service ./
COPY --from=group-builder /app/config/local.yaml ./
ENTRYPOINT ["./group-service"]
CMD ["--config=./local.yaml"]