FROM golang:latest
WORKDIR "/home/apps/credit_cards"
COPY . /app
RUN go mod download \
    && go mod verify \

RUN go build -o companies -a .
EXPOSE 3001
CMD ["./companies"]