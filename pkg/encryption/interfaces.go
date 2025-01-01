package encryption

type Encryptor interface {
	Encrypt(data Data) (string, error)
	Decrypt(data Data) (string, error)
	EncryptStruct(data interface{}, key string) (interface{}, error)
	DecryptStruct(data interface{}, key string) (interface{}, error)
}
