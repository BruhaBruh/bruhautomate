#!/usr/bin/env bash

go build -o bruhautomate .

sudo rm -f /usr/local/bin/bruhautomate
sudo rm -f /usr/local/bin/bam
sudo cp bruhautomate /usr/local/bin/bruhautomate
sudo ln -s /usr/local/bin/bruhautomate /usr/local/bin/bam
