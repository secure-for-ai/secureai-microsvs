db.createUser({user:"test", pwd: "password",  roles: [{role: "readWrite", db: "gtest" }], mechanisms:["SCRAM-SHA-1"],passwordDigestor:"server" })
db.createCollection("user")
