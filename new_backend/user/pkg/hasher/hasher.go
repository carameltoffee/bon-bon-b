package hasher

type Hasher interface {
	Hash(password string) (string, error)
	Compare(password, hashed string) bool
}
