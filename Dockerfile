FROM golang:1.26 AS builder
WORKDIR /src

COPY ./go.sum ./go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /gen -ldflags="-w -s"
RUN /gen -out /out -base-url https://sh.devminer.xyz


FROM ghcr.io/thedevminertv/gostatic:1.5.2

COPY --from=builder /out /static