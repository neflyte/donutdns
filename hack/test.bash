#!/usr/bin/env bash
#
# test.bash -- Run a donutdns server with testing configuration values
#
export DONUT_DNS_NO_DEFAULTS="true"
export DONUT_DNS_NO_LOG="false"
export DONUT_DNS_NO_DEBUG="false"
export DONUT_DNS_PORT=15353
export DONUT_DNS_ALLOW_FILE="hack/allow_file.txt"
export DONUT_DNS_ALLOWSUFFIX_FILE="hack/allowsuffix_file.txt"
export DONUT_DNS_BLOCK_FILE="hack/block_file.txt"
export DONUT_DNS_UPSTREAM_1="192.168.1.1"
export DONUT_DNS_UPSTREAM_MAX_FAILS="1"
./donutdns
