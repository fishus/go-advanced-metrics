package grpc

type Config struct {
	ServerAddr string
	SecretKey  string
	PublicKey  []byte
}
