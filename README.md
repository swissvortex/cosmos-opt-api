# cosmos-opt-api
Enhance prometheus metrics using the Cosmos API endpoints


## How to:
```
git clone https://github.com/swissvortex/cosmos-opt-api.git

cd cosmos-opt-api/

go build

./cosmos-opt-api
```

## Create systemd service
```
sudo bash -c 'echo "[Unit]
Description=Cosmos Opt API
After=network-online.target
[Service]
User=user
ExecStart=/home/user/cosmos-opt-api/cosmos-opt-api
Restart=on-failure
RestartSec=3
[Install]
WantedBy=multi-user.target" > /etc/systemd/system/cosmos-opt-api.service'
```

`sudo systemctl start cosmos-opt-api.service && sudo systemctl enable cosmos-opt-api.service`

> Check logs with `journalctl -f -u cosmos-opt-api.service`