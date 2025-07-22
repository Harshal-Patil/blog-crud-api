#!/bin/sh
chmod go-w filebeat.yml
exec filebeat -e -c filebeat.yml