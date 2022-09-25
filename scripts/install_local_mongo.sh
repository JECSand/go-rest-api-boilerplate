#!/bin/bash

if [ "$OSName" = "CentOS Linux" ]
 then
    echo '[MongoDB]
name=MongoDB Repository
baseurl=http://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/4.0/x86_64/
gpgcheck=0
enabled=1' | sudo tee -a /etc/yum.repos.d/mongodb.repo
    sudo yum install -y mongodb-org
    sudo systemctl start mongod.service
    sudo systemctl enable mongod.service
else
    sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4
    echo 'deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.0 multiverse' | sudo tee /etc/apt/sources.list.d/mongodb.list
    sudo apt -y update && sudo apt -y install mongodb-org
    sudo systemctl enable mongod
    sudo systemctl start mongod
fi