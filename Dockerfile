    FROM ubuntu:20.04
    
    WORKDIR /libsrt
    ENV DEBIAN_FRONTEND=noninteractive
    
    RUN apt-get update && apt-get install -y git build-essential
    RUN apt-get install -y tclsh pkg-config cmake libssl-dev curl
    
    # INSTALL libsrt
    RUN git clone https://github.com/Haivision/srt.git /libsrt
    RUN ./configure
    RUN make
    RUN make install
    

    WORKDIR /app
    # setup go environment
    RUN curl -LO https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
    RUN tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

    ENV PATH=$PATH:/usr/local/go/bin
    ENV GOPATH=/root/go
    ENV PATH=$PATH:$GOPATH/bin

    RUN go version

    COPY . .

    # BUILD srt-gateway
    RUN go build -buildvcs=false -o srt-gateway  . 
    
    #COPY ./srt-gateway /app/
    
    EXPOSE 9800
    EXPOSE 9801
    EXPOSE 9802
    

    ENV LD_LIBRARY_PATH=/usr/local/lib

    RUN chmod +x /app/srt-gateway
    
    ENTRYPOINT ["/app/srt-gateway"]