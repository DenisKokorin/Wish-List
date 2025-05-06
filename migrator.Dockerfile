FROM alpine/git AS clone

WORKDIR /

RUN git clone https://github.com/Gergenus/AuthService.git ./AuthService


FROM golang:1.24.2 AS migrator-builder

WORKDIR /app

COPY . .

COPY ./migrations cmd/migrator/migrations

COPY --from=clone ./AuthService/internal/migrations cmd/migrator/migrations

RUN cd cmd/migrator && CGO_ENABLED=0 GOOS=linux go build -o migrator main.go


FROM alpine:latest

WORKDIR /

COPY --from=migrator-builder /app/cmd/migrator ./

ENTRYPOINT ["./migrator"]
