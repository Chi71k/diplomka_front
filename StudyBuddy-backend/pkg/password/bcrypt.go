package password

import "golang.org/x/crypto/bcrypt"

const cost = bcrypt.DefaultCost

// Hash returns a bcrypt hash of the password.
func Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare returns true if plain matches the hash.
func Compare(hash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
