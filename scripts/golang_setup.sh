#!/bin/bash

set -ex

# Install go (version 1.23.0) and add to local path
sudo rm -rf /usr/local/go
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.23.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source $HOME/.bashrc
rm go1.23.0.linux-amd64.tar.gz

set +ex
