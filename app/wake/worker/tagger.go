package worker

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/coralproject/pillar/pkg/model"
)

const (
	ADD_USER_MARKER string    = "add_user_markers"
	REMOVE_USER_MARKER string = "remove_user_markers"
)

func UpdateUserTag(event model.Event) {

	//Call an internal End-Point

}

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getSubject(event string) string {
	switch event {
	case model.EventTagAdded:
		return ADD_USER_MARKER

	case model.EventTagRemoved:
		return REMOVE_USER_MARKER
	}

	return ""
}
