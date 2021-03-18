package auth

// TODO: I have not auditted this google-auth-id-token-verifier lib yet!

import (
	"fmt"

	verifier "github.com/futurenda/google-auth-id-token-verifier"
)

var v = verifier.Verifier{}
var aud = "656765689994-f29hh63in3j1362mom3ek00ukcmru8jq.apps.googleusercontent.com"

func googleOAuthTokenParser(token string) (isValid bool, unique_user_id string) {
	err := v.VerifyIDToken(token, []string{
		aud,
	})
	if err != nil {
		fmt.Println("error verity id token google login", err.Error())
		return false, ""
	}
	claimSet, err := verifier.Decode(token)
	if err != nil {
		fmt.Println("error id token verifier google login", err.Error())
		return false, ""
	}
	return true, "__google__" + claimSet.Sub
}
