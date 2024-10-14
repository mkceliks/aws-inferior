FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY lambda/ .

RUN GOOS=linux GOARCH=amd64 go build -o main

FROM public.ecr.aws/lambda/go:latest

COPY --from=builder /app/main ${LAMBDA_TASK_ROOT}

CMD ["master"]
