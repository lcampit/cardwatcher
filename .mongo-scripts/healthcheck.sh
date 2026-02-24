#! /bin/bash

mongosh -u root -p "${MONGO_INITDB_ROOT_PASSWORD}" --authenticationDatabase admin --quiet \
  --eval "db.runCommand({ping:1}).ok" | grep 1 >/dev/null
