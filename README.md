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