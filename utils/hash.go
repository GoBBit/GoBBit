package utils

import (
	"encoding/base64"
	"crypto/sha512"
	//"gopkg.in/mgo.v2/bson"
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
func GenerateUserSession(uid, site_key string)(string){
	// Generate usersession based on UID, timestamp and hash of:
	// 	UID+Timestamp+SITE_KEY
	// Example: UID:Timestamp:HASH
	now := time.Now().Unix() * 1000
	nowStr := strconv.FormatInt(now, 10)
	hash := CalculateHash(uid + nowStr + site_key)
	session := uid + ":" + nowStr + ":" + hash

	return session
}

func CheckSession(uid, timestamp, site_key, hash string)(bool){
	// Check the cookie
	calculatedHash := CalculateHash(uid + timestamp + site_key)

	return (calculatedHash == hash)
}

