#!/usr/bin/env bash

sudo apt update && sudo apt upgrade
sudo apt install autoconf automake libtool curl make g++ unzip -y
mkdir ~/tmp
cd ~/tmp
git clone https://github.com/google/protobuf.git
cd protobuf
git submodule update --init --recursive
./autogen.sh
./configure
make
make check
sudo make install
sudo ldconfig
