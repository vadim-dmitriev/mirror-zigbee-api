# docker run -d --restart=unless-stopped -name mirror-zigbee-api <imageID>

FROM golang:1.19-alpine

WORKDIR $GOPATH/src/github.com/vadim-dmitriev/mirror-zigbee-api

COPY go.mod ./
COPY go.sum ./

COPY . .

RUN go build -o mirror-zigbee-api cmd/main.go

CMD ["./mirror-zigbee-api"]
