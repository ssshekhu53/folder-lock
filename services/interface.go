package services

type FolderLock interface {
	Init(password string) error
	Lock() error
	Unlock(password string) error
	UpdatePassword(oldPassword, newPassword string) error
}

type Crypt interface {
	Encrypt(creds []byte) []byte
	Decrypt(cred string) ([]byte, error)
}
