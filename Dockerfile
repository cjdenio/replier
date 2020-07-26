FROM golang:latest

WORKDIR /usr/src/app

COPY . .

RUN go get .

RUN go get -u github.com/cosmtrek/air
ENV air_wd /usr/src/app

EXPOSE 3000

CMD [ "air" ]