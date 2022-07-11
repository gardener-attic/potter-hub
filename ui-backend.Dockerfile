FROM eu.gcr.io/gardener-project/3rd/golang:1.17.9 as builder

WORKDIR /app

RUN apt-get install git
RUN go get github.com/go-delve/delve/cmd/dlv

COPY go.mod go.sum ./
COPY pkg/ pkg/
COPY cmd/ui-backend/ ./cmd/ui-backend
COPY cmd/apprepository-controller/ cmd/apprepository-controller/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o ui-backend cmd/ui-backend/main.go

#### BASE ####
FROM gcr.io/distroless/static-debian11:nonroot AS base

WORKDIR /app

COPY --from=builder app/ui-backend /
COPY --from=builder /go/bin/dlv /dlv

CMD ["/ui-backend"]