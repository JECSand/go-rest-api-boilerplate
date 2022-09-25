#!/bin/bash
if [ "$OSName" = "CentOS Linux" ]
 then
    echo "[Unit]
Description=rest api initialization
[Service]
PIDFile=/tmp/restapi.pid-4040
User="$2"
Group="$2"
WorkingDirectory="$1"
ExecStart=/bin/bash -c '"$1"/go-rest-api-boilerplate'
[Install]
WantedBy=multi-user.target" >> /usr/lib/systemd/system/restapi.service
else
  echo "[Unit]
Description=golang rest api initialization
[Service]
PIDFile=/tmp/restapi.pid-4040
User="$2"
Group="$2"
WorkingDirectory="$1"
ExecStart=/bin/bash -c '"$1"/go-rest-api-boilerplate'
[Install]
WantedBy=multi-user.target" >> /lib/systemd/system/restapi.service
fi