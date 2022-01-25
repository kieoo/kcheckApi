#!/bin/sh
pkill kcheckApi
export GOGIN_MODE=release
export GOPORT=8001
nohup ./bin/kcheckApi & 2>&1