FROM golang:1.22 AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /src

COPY . .

RUN apt-get install -yf git

RUN go build \
  -trimpath \
  -installsuffix cgo \
  -ldflags "-extldflags '-static'" \
  -buildvcs=true \
  -o /bin/wtgbot \
  ./cmd/warehouse-bot

FROM golang:alpine
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add curl
COPY --from=builder /bin/wtgbot /bin/wtgbot

ENTRYPOINT ["/bin/wtgbot"]
