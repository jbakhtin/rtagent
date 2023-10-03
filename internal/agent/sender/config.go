package sender

type Configer interface {
	GetServerAddress() string
	GetKeyApp() string
	GetCryptoKey() string
}
