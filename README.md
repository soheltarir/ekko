# Ekko

![GitHub release (latest by date)](https://img.shields.io/github/v/release/soheltarir/ekko?style=flat-square)
![GitHub](https://img.shields.io/github/license/soheltarir/ekko?style=flat-square)

Tool to ping servers (using ICMP & UDP) and track response time and packet loss.

![Sample](sample_ui.png)

## 📦 Installation
Download the appropriate zip file from [Ekko Releases](https://github.com/soheltarir/ekko/releases) link. Extract the 
zip file in a directory location of your choosing.

## 🛠 Configuration

The Ekko binary reads configuration from the `config.yaml` file which should reside in the same folder as that of the
binary/executable. By default, the release package comes pre-configured with a `config.yaml` which you could change as per your
requirements.

## 🪀 Usage

### MacOS

1. [Open a terminal](https://support.apple.com/en-in/guide/terminal/apd5265185d-f365-44cb-8b09-71a064a42125/mac) and go 
the directory where you extracted the zip file; `cd <YourInstallDirectory>/ekko`
2. Ekko requires super admin privileges to run, hence run the program using `sudo ./ekko`
3. Press `control` and `C` to close the program or quit the terminal window.

### Linux
1. Open a terminal and go the directory where you extracted the zip file; `cd <YourInstallDirectory>/ekko`
2. Ekko requires super admin privileges to run, hence run the program using `sudo ./ekko`
3. Press `Ctrl` and `C` to close the program or quit the terminal window.

### Windows
1. Go the folder where the zip file is extracted using Windows Explorer.
2. Right click `ekko.exe`, and select "Run as Administrator".

## Logs
Network statistics UI and file logs are enabled by default (which you could configure on your in the `config.yaml` file).
The file logs reside in the `logs` folder in the directory wherein the package is extracted. You would find logs files
generated therein; **results.ndjson** stores network statistics for pings generated and **debug.ndjson** stores 
application debug logs.
```json lines
{"severity":"info","timestamp":"2022-01-09T19:39:21.186+0530","message":"Ekko service started"}
{"severity":"info","timestamp":"2022-01-09T19:39:21.231+0530","message":"Ping started","worker_id":0,"server_name":"Valorant (Mumbai 1)","server_ip":"75.2.66.166","labels":{"game":"Valorant","provider":"Riot"}}
{"severity":"info","timestamp":"2022-01-09T19:39:21.231+0530","message":"Ping started","worker_id":4,"server_name":"Valorant (Behrain 1)","server_ip":"75.2.105.73","labels":{"game":"Valorant","provider":"Riot"}}
{"severity":"info","timestamp":"2022-01-09T19:39:21.231+0530","message":"Ping started","worker_id":2,"server_name":"Valorant (Mumbai 2)","server_ip":"99.83.136.104","labels":{"game":"Valorant","provider":"Riot"}}
{"severity":"info","timestamp":"2022-01-09T19:39:21.231+0530","message":"Ping started","worker_id":3,"server_name":"Valorant (Behrain 2)","server_ip":"99.83.199.240","labels":{"game":"Valorant","provider":"Riot"}}
{"severity":"error","timestamp":"2022-01-09T19:39:21.234+0530","message":"Failed to initialise ping","worker_id":1,"server_name":"Dota2 (SEA-1)","server_ip":"sgp-1.valve.net","labels":{"game":"Dota2","provider":"Valve"},"error":"lookup sgp-1.valve.net: no such host"}
{"severity":"info","timestamp":"2022-01-09T19:39:21.272+0530","message":"Ping started","worker_id":1,"server_name":"Dota2 (SEA-2)","server_ip":"sgp-2.valve.net","labels":{"game":"Dota2","provider":"Valve"}}
{"severity":"info","timestamp":"2022-01-09T19:39:26.276+0530","message":"Ping complete","worker_id":2,"server_name":"Valorant (Mumbai 2)","server_ip":"99.83.136.104","labels":{"game":"Valorant","provider":"Riot"},"num_packets":6,"packet_loss":0,"avg_rtt":66,"min_rtt":40,"max_rtt":112}
{"severity":"info","timestamp":"2022-01-09T19:39:31.611+0530","message":"Ping complete","worker_id":1,"server_name":"Dota2 (SEA-2)","server_ip":"sgp-2.valve.net","labels":{"game":"Dota2","provider":"Valve"},"num_packets":11,"packet_loss":0,"avg_rtt":102,"min_rtt":60,"max_rtt":247}
{"severity":"info","timestamp":"2022-01-09T19:39:35.303+0530","message":"Ping complete","worker_id":3,"server_name":"Valorant (Behrain 2)","server_ip":"99.83.199.240","labels":{"game":"Valorant","provider":"Riot"},"num_packets":15,"packet_loss":0,"avg_rtt":106,"min_rtt":41,"max_rtt":276}
{"severity":"info","timestamp":"2022-01-09T19:39:35.373+0530","message":"Ping complete","worker_id":4,"server_name":"Valorant (Behrain 1)","server_ip":"75.2.105.73","labels":{"game":"Valorant","provider":"Riot"},"num_packets":15,"packet_loss":0,"avg_rtt":116,"min_rtt":49,"max_rtt":278}
{"severity":"info","timestamp":"2022-01-09T19:39:38.368+0530","message":"Ping complete","worker_id":0,"server_name":"Valorant (Mumbai 1)","server_ip":"75.2.66.166","labels":{"game":"Valorant","provider":"Riot"},"num_packets":18,"packet_loss":0,"avg_rtt":98,"min_rtt":44,"max_rtt":278}
```
The log output is in newline-delimited JSON format (learn more here: http://ndjson.org/), upon which you could generate
metrics later on to trigger alerts or create historical dashboards to track network performance of the destinations configured.
