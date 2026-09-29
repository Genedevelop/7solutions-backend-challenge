package outbound

type IPasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
