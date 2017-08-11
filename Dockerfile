FROM golang:onbuild
RUN mkdir /paxos
ADD . /paxos/
WORKDIR /paxos
RUN go build -o main .
CMD ["/paxos/main"]
EXPOSE 5000


