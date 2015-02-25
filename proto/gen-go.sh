#!/bin/bash

set -eux

if [[ ! -d SteamKit ]]; then
	git clone https://github.com/SteamRE/SteamKit.git
fi
SRC_DIR=SteamKit/Resources/Protobufs

SRCS="steammessages.proto gcsystemmsgs.proto base_gcmessages.proto gcsdk_gcmessages.proto econ_gcmessages.proto steammessages_cloud.steamworkssdk.proto steammessages_oauth.steamworkssdk.proto steammessages_publishedfile.steamworkssdk.proto network_connection.proto dota_gcmessages_common.proto dota_gcmessages_client.proto dota_gcmessages_client_fantasy.proto
dota_gcmessages_server.proto"

for x in $SRCS; do
	protoc --go_out=. -I $SRC_DIR -I $SRC_DIR/dota $SRC_DIR/dota/$x
done
