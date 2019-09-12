FROM golang:1.12
ENV GO111MODULE=on
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/

RUN apt-get -qq update \
  && apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    tesseract-ocr
RUN tesseract -v
# Load best data
# RUN rm /usr/share/tesseract-ocr/4.00/tessdata/eng.traineddata
# RUN wget -P /usr/share/tesseract-ocr/4.00/tessdata/ \
#   "https://github.com/tesseract-ocr/tessdata_best/blob/master/eng.traineddata"

ADD . $GOPATH/src/github.com/jou66jou/go-orcnums-xy
WORKDIR $GOPATH/src/github.com/jou66jou/go-orcnums-xy
RUN go get github.com/otiai10/...
# RUN go test -v github.com/otiai10/gosseract
RUN GOOS=linux GOARCH=amd64 go build -o main .

CMD ./main -p 8080
# CMD tail -f /dev/null