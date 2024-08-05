FROM alpine:latest

COPY server ./server

WORKDIR ./server

#delete test in order to not install testify
RUN rm -rf internal/test

#install golang version 1.22.5
RUN wget https://golang.org/dl/go1.22.5.linux-amd64.tar.gz &&  \
    tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz && \
    rm -rf go1.22.5.linux-amd64.tar.gz

ENV GOPATH=/root/go
ENV PATH=${GOPATH}/bin:/usr/local/go/bin:$PATH
ENV GOBIN=$GOROOT/bin

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

ENV GO111MODULE=on

RUN go mod tidy && go build cmd/main/main.go

EXPOSE 18080

ENTRYPOINT ["./main"]
