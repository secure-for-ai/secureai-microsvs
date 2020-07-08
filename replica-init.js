db.createUser({user:"test", pwd: "password",  roles: [{role: "readWrite", db: "gtest" }], mechanisms:["SCRAM-SHA-1"],passwordDigestor:"server" })
db.createCollection("user", { collation: { locale: 'en_US', strength: 2 } } )
db.user.createIndex( { "username": 1 }, { unique: true, collation: { locale: 'en_US', strength: 2 } } )
