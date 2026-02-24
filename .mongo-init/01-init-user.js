// Create application user on first startup

db = db.getSiblingDB(process.env.MONGO_DATABASE);

db.createUser({
  user: process.env.MONGO_USERNAME,
  pwd: process.env.MONGO_PASSWORD,
  roles: [
    { role: "readWrite", db: process.env.MONGO_DATABASE }
  ],
  mechanisms: ["SCRAM-SHA-256"]
});

