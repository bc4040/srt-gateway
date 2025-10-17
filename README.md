 # srt-gateway
 This project provides a gateway service for SRT live video transmission and is being used as a 24/7 live production confidence monitoring solution.
 
 ### An example use case:
 SDI to SRT hardware encoders on a local LAN can be exposed to the internet via a secure "double NIC'ed" host running srt-gateway.

 The hardware encoders connect to srt-gateway via SRT caller mode.
 srt-gateway makes this feed available on the two sender channels in SRT listener mode.
 Internet clients connect to srt-gateway and receive the SRT stream on the sender channels in caller mode.


 srt-gateway provides one ingest channel and two sender channels.

 This can be multiplied at scale by running multiple instances in Docker.
 
 # Running via Docker
 Although it can run as a standalone binary, the project is intended to be run inside a Docker container.

 There are 3 passphrases you can define.  The first two are mandatory.
 ```
 PASSPHRASE_IN - The passphrase for the ingest port (required)
 PASSPHRASE_OUT1 - The passphrase of the first sender (required)
 PASSPHRASE_OUT2 - The passphrase of the second sender (when defined, second sender is enabled)
 ```
 

 `docker run --rm -it -e "PASSPHRASE_IN=password" -e "PASSPHRASE_OUT1=password" srt-gateway`

 There are 3 ports you can map.
 These are hardcoded in the binary, but can be modified via your Docker container configuration.
 ```
 :9800 - The ingest port (input) (srt listen mode)
 :9801 - The first sender (srt listen mode)
 :9802 - The second sender (srt listen mode)
 ```

We can further define specific port mapping with Docker

Example:
We want to use port 5000 on the local LAN for ingest.  In this example the LAN IP address of the srt-gateway host is 192.168.1.100.

We want to use port 5050 (sender 1) and 5051 (sender 2) on the WAN.  We'll define these on all interfaces.

 `docker run --rm -it -e "PASSPHRASE_IN=password" -e "PASSPHRASE_OUT1=password" -e "PASSPHRASE_OUT2=password"  -p 192.168.1.100:5000:9800 -p 5050:9801 -p 5051:9802 srt-gateway`

Now you can point your hardware encoder to 192.168.1.100:5000 in caller mode, and it will be relayed to WAN ports 5050 and 5051.  

Connect to the srt-gateway host with a client viewer in caller mode to srt-gateway on either port 5050 or 5051.