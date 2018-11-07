FROM golang:1.9.4

ADD . /go/src/virtualizer
WORKDIR /go/src/virtualizer

RUN go get github.com/BurntSushi/toml
RUN go get github.com/Sirupsen/logrus
RUN go get github.com/clbanning/mxj
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/tidwall/gjson

RUN go build main/main.go

ENTRYPOINT ["./main"]
