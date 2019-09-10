FROM golang:1.12
ENV GO111MODULE=on


RUN apt-get -qq update \
  && apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    tesseract-ocr

# Load languages
RUN apt-get install -y \
  tesseract-ocr-jpn

ADD . $GOPATH/src/github.com/jou66jou/go-orcnums-xy
WORKDIR $GOPATH/src/github.com/jou66jou/go-orcnums-xy
RUN go get ./...
RUN go test -v github.com/otiai10/gosseract

CMD $GOPATH/bin/ocrserver -p 8080
