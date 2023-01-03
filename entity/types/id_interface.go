package types

type IdInterface interface {
	~string | ~int
	String() string
}

