FROM dogfooter/golang:1.11.5-dev as build
ENV DOGFOOTER_HOME /go/src
RUN git clone https://wonsuck_song:wssong98@gitlab.lazybird.kr/illuco/dogfooter-image.git
RUN git clone https://wonsuck_song:wssong98@gitlab.lazybird.kr/illuco/dogfooter-patient.git
RUN git clone https://wonsuck_song:wssong98@gitlab.lazybird.kr/illuco/dogfooter-date.git
RUN git clone https://wonsuck_song:wssong98@gitlab.lazybird.kr/illuco/dogfooter-diagnosis.git
WORKDIR /go/src/dogfooter-control
ADD . .
RUN go build -o app_dogfooter_control control/cmd/main.go

FROM node:8.15.1 as frontend_build
WORKDIR /usr/src
RUN git clone https://wonsuck_song:wssong98@gitlab.lazybird.kr/illuco/dogfooter.git app
WORKDIR /usr/src/app
RUN echo "\nNODE_ENV=production\n" >> .env
RUN yarn
ENV PATH /usr/src/app/node_modules/.bin:$PATH
RUN yarn build

FROM alpine:3.9
ENV DOGFOOTER_HOME /var/local
ENV DOGFOOTER_DATA /var/local/data
ENV DOGFOOTER_FRONTEND_HOME frontend
WORKDIR /var/local/dogfooter-control/config
COPY --from=build /go/src/dogfooter-control/config .
WORKDIR /var/local/dogfooter-control/img
COPY --from=build /go/src/dogfooter-control/img .
WORKDIR /usr/local/bin/frontend
COPY --from=frontend_build /usr/src/app/build .
WORKDIR /usr/local/bin
COPY --from=build /go/src/dogfooter-control/app_dogfooter_control /usr/local/bin/app_dogfooter_control

ENTRYPOINT ["app_dogfooter_control"]
