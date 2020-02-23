FROM golang

WORKDIR /app/news-consumer

COPY ./go.mod ./

RUN go mod download

COPY . .

RUN go test ./...

RUN go build -o ./bin/news-consumer .

CMD [ "./bin/news-consumer" ]