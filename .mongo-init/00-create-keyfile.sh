#!/bin/bash
# Create a keyfile to use in the replica set if it is missing

if [ ! -f /data/configdb/mongo-keyfile ]; then
  openssl rand -base64 756 >/data/configdb/mongo-keyfile
  chmod 600 /data/configdb/mongo-keyfile
fi
