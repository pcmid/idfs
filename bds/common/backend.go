package common

type OssBackend interface {
	Get(item string) ([]byte, error)
	Put(item string, data []byte)
	Delete(item string)
}
