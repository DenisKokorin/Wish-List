FROM golang:1.24.2 AS wishlist-builder
WORKDIR /app
COPY . .
RUN cd cmd/wishlist && CGO_ENABLED=0 GOOS=linux go build -o wishlist-service main.go

FROM alpine AS wishlist
WORKDIR /
COPY --from=wishlist-builder /app/cmd/wishlist/wishlist-service ./
COPY --from=wishlist-builder /app/config/local.yaml ./
ENTRYPOINT ["./wishlist-service"]
CMD ["--config=./local.yaml"]
