package secret

type Secret struct {
	Name            string
	Data            map[string][]byte
	ResourceVersion string
}

type SecretInterface interface {
	V1() V1
}

type V1 interface {
	Create(*Secret) (*Secret, error)
	Update(*Secret) (*Secret, error)
	Delete(string) error
	Get(string) (*Secret, error)
	List() ([]*Secret, error)
}
