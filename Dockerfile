FROM golang:1.12.6

ENV TINI_VERSION v0.18.0

ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

WORKDIR /workspace/knaing-receiver
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /app/stan-to-logz .
RUN rm -rf $GOPATH/pkg/mod

ENTRYPOINT ["/tini", "-s", "--"]
CMD ["/app/stan-to-logz"]