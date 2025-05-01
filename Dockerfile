FROM golang:1.23 AS wishlist-builder
WORKDIR /app
COPY . .
COPY .env cmd/wishlist
RUN cd cmd/wishlist && CGO_ENABLED=0 GOOS=linux go build -o wishlist-service main.go

#FROM alpine:latest

#FROM alpine AS wishlist
#COPY --from=wishlist-builder /wishlist-service /app/
#COPY --from=wishlist-builder /config/local.yaml /app/
CMD ["./cmd/wishlist/wishlist-service", "--config=./config/local.yaml"]
#ENTRYPOINT ["/app/wishlist-service"]
#CMD ["--config=/app/local.yaml"]