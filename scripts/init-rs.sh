#! /bin/bash
mongosh --host mongo -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" --authenticationDatabase \
  admin --quiet --eval '
        try {
          const s = rs.status();
          if (s.ok === 1) { print("Replica set already initialized"); quit(0); }
        } catch (e) {
          if (e.codeName === "NotYetInitialized" || e.code === 94) {
            print("Initializing replica set rs0...");
            rs.initiate({_id:"rs0", members:[{_id:0, host:"mongo:27017"}]});
            // wait until PRIMARY
            for (let i = 0; i < 30; i++) {
              try {
                const st = rs.status();
                if (st.members && st.members.some(m => m.stateStr === "PRIMARY")) { print("Replica set PRIMARY"); quit(0); }
              } catch (_) {}
              sleep(1000);
            }
            print("Replica set did not reach PRIMARY in time"); quit(1);
          }
          print("Unexpected error: " + tojson(e)); quit(1);
        }
      '
