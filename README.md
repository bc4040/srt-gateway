 # srt-gateway
 This project provides a secure gateway service for SRT live video transmission and is currently being used as a 24/7 live production confidence monitoring solution.
 
 ### An example use case:
 You have a SDI/HDMI to SRT hardware encoder on a local LAN and want to securely expose it to the internet. 
 srt-gateway can sit between the encoder (LAN) and the internet (WAN).  Remote viewers connect to srt-gateway to view the encoded stream.

 This prevents the need to expose an insecure hardware encoder appliance to the internet.

 srt-gateway provides one ingest (input) channel and two sender (output) channels.   Currently, only one client can connect to a single channel at a time.

 This can be multiplied at scale by running multiple instances in Docker.
 
## Build for Docker

Clone this repository then run:
`docker build -t srt-gateway:latest .`

 ## Running via Docker
 Although it can run as a standalone binary, the project is intended to be run inside a Docker container.

 There are 3 passphrases you can define.  The first two are mandatory.
 ```
 PASSPHRASE_IN - The passphrase for the ingest port (required)
 PASSPHRASE_OUT1 - The passphrase of the first sender (required)
 PASSPHRASE_OUT2 - The passphrase of the second sender (when defined, second sender is enabled)
 ```
 
 `docker run --rm -it -e "PASSPHRASE_IN=password" -e "PASSPHRASE_OUT1=password" srt-gateway`

### Custom port mapping

 Ports are hardcoded in the binary, but can be exposed to the outside world on different ports via your Docker container configuration.
 There are 3 ports you can map.
 ```
 :9800 - The ingest port (input) (srt listen mode)
 :9801 - The first sender (srt listen mode)
 :9802 - The second sender (srt listen mode)
 ```
Note: From the srt-gateway's perspective, all SRT connections are in listener mode.  
This means both the encoder and clients will connect to srt-gateway in *caller* mode.

We can further define specific port mapping with Docker...

### Example:
We want to use port 5000 on the local LAN for ingest.  
In this example the LAN IP address of the srt-gateway host is 192.168.1.100.

We want to use port 5050 (sender 1) and 5051 (sender 2) on the WAN.  We'll define these on all interfaces.

 `docker run --rm -it -e "PASSPHRASE_IN=password" -e "PASSPHRASE_OUT1=password" -e "PASSPHRASE_OUT2=password"  -p 192.168.1.100:5000:9800 -p 5050:9801 -p 5051:9802 srt-gateway`


Now you can point your hardware encoder to 192.168.1.100:5000 in caller mode, and it will be relayed to WAN ports 5050 and 5051.  

Connect to the srt-gateway host with a client viewer in caller mode to srt-gateway on either port 5050 or 5051.
