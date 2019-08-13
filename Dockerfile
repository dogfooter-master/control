FROM dogfooter/golang:1.11.5-dev as build
WORKDIR /go/src/dogfooter-control
ADD . .
#RUN apk add --no-cache bash git openssh
#RUN dep init -v -no-examples
RUN go build -o app_dogfooter_control dogfooter-control/control/cmd

FROM alpine:3.9
ENV DOGFOOTER_HOME /var/local
WORKDIR /var/local/dogfooter-control/config
COPY --from=build /go/src/dogfooter-control/config .
WORKDIR /var/local/dogfooter-control/img
COPY --from=build /go/src/dogfooter-control/img .
WORKDIR /usr/local/bin
COPY --from=build /go/src/dogfooter-control/app_dogfooter_control /usr/local/bin/app_dogfooter_control

ENTRYPOINT ["app_dogfooter_control"]
