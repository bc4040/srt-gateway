# Overview
This is a simple implementation of an SRT (Secure Reliable Transport) gateway that allows you to expose an SRT encoder or stream from your private LAN to the public internet via a dedicated relay host.

Right now it supports replicating a single ingest stream to two output channels, and only one connection is allowed per channel.  All connections to/from the srt-gateway should be in SRT "caller" mode

The ports for each channel are currently hard-coded.
- Ingest: 9800/udp
- Sender 1: 9801/udp
- Sender 2: 9802/udp

**It works like this:**

- Point your encoder to the ingest channel of srt-gateway
- Connect to one of the sender channels via an SRT player such as VLC, Larix, VMIX, etc. to view the stream

# How to Run

 You will need to define some environment variables at runtime
 
 That includes passphrases for the ingest (IN) channel as well as the sender (OUT) channels.

 Optionally, you can elect to record the stream with the env var "record" set to "true"
 A .ts file will be created in the same directory as the binary.

```
PASSPHRASE_IN=
PASSPHRASE_OUT1=
PASSPHRASE_OUT2=
record=true
```

# Build with Docker

```
docker build -t srt-gateway .
```

# Run with Docker

Be sure to set the environment variables in your Docker run command

 
 ```
 docker run --rm -it -e "PASSPHRASE_IN=<passwordhere>" -e "PASSPHRASE_OUT1=<passwordhere>" srt-gateway
```