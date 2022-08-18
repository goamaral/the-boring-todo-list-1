package m3o

type provider struct {
}

type Provider interface {
}

func NewM3OProvider() *provider {
	return &provider{}
}
