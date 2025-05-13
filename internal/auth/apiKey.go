package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {

	auth, ok := headers["Authorization"]
	// fmt.Println(auth)
	if ok {
		if len(auth) == 1{
			return strings.Replace(auth[0], "ApiKey ", "", 1), nil
		}
		return "", fmt.Errorf("something is wrong with this header, %v", auth)
	}

	return "", errors.New("couldn't find the right header")
}