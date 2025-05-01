package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestJWTCreationAndValidation(t *testing.T) {

	uid := uuid.New()
	secret := "beans"

	token, err := MakeJWT(uid, secret)

	if err != nil {
		t.Errorf("issue with JWT creation: %s\n", err.Error())
	}

	checkId, errV := ValidateJWT(token, secret)
	if errV != nil {
		t.Errorf("issue with JWT validation: %s\n", errV.Error())
	}
	if uid != checkId {
		t.Errorf("these uuids are not the same: %s | %s", uid.String(), checkId.String())
	}

}