package utils

import (
	"encoding/base64"
	"crypto/sha512"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strconv"
)

// Hash Functions
func CalculateHash(s string) (string){
    sha512.New()
    sBytes := []byte(s)
    hash := sha512.Sum512(sBytes)
    b64hash := base64.URLEncoding.EncodeToString(hash[:])

    return b64hash
}

// Session Function
func GenerateUserSession()(string){
	// Generate usersession based on an mongodb objectID and the actual timestamp
	now := time.Now().Unix() * 1000
	nowStr := strconv.FormatInt(now, 10)
	id := bson.NewObjectId()
	session := nowStr + id.Hex()

	return CalculateHash(session)
}

