db = db.getSiblingDB("octocv")
db.createUser({user: "octocv", pwd: "root", roles: [{role: "readWrite", db: "octocv"}]})