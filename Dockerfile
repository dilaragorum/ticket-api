FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CG0_ENABLED=0 go build -o api

FROM alpine

COPY --from=0 /app/api ./api
COPY --from=0 /app/.env.dev ./.env.dev

EXPOSE 3000

CMD [ "./api" ]