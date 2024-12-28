Example Dockerfile:
    FROM ubuntu:20.04
    
    WORKDIR /app
    ENV DEBIAN_FRONTEND=noninteractive
    
    RUN apt-get update && apt-get install -y git build-essential
    
    RUN apt-get install -y tclsh pkg-config cmake libssl-dev
    
    RUN git clone https://github.com/Haivision/srt.git /app
    
    RUN ./configure
    RUN make
    
    RUN make install
    
    
    ENV LD_LIBRARY_PATH=/usr/local/lib
    
    COPY ./srt-gateway /app/
    
    
    EXPOSE 9800
    EXPOSE 5000
    EXPOSE 5009
    
    RUN chmod +x /app/srt-gateway
    
    ENTRYPOINT ["/app/srt-gateway"]