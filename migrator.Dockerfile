FROM golang:1.23
WORKDIR /app
COPY . .
COPY .env cmd/migrator
COPY migrations cmd/migrator
RUN cd cmd/migrator && go build -o migrator main.go
CMD ["./cmd/migrator/migrator"]