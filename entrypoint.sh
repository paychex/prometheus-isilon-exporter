#!/bin/sh
# B"H

# Since CMD instruction can't get envvironment variables moved the start command here.
/app/prometheus-isilon-exporter --username ${USERNAME} --password ${PASSWORD} --url ${URL}