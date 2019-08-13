FROM dogfooterm/golang:1.11.5-dev
ENV DOGFOOTER_HOME /go/src
WORKDIR /go/src/dogfooter-control
ADD . .
#RUN dep init -v -no-examples
#RUN dep ensure -v -vendor-only
#RUN go get -v dogfooter-control/control/cmd
RUN chmod +x Build.sh
ENTRYPOINT sh Build.sh
