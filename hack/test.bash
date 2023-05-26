#!/usr/bin/env bash
#
# test.bash
#
export DONUT_DNS_NO_DEFAULTS="true"
export DONUT_DNS_NO_LOG="false"
export DONUT_DNS_NO_DEBUG="false"
export DONUT_DNS_PORT=15353
export DONUT_DNS_ALLOW_FILE="hack/allow_list.txt"
export DONUT_DNS_ALLOWSUFFIX_FILE="hack/allowsuffix.txt"
export DONUT_DNS_BLOCK_FILE="hack/block_list.txt"
export DONUT_DNS_UPSTREAM_1="1.1.1.1"
export DONUT_DNS_UPSTREAM_2="1.0.0.1"
output/donutdns
