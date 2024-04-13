package windows

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"ssshekhu53/folder-lock/services"
)

const (
	fileName              = `private`
	encryptedDataFileName = `.encrypted-data`
)

type windows struct {
	crypt services.Crypt
}

func New(crypt services.Crypt) services.FolderLock {
	return &windows{crypt: crypt}
}

func (m *windows) Init(password string) error {
	_, err := os.Stat(fileName)
	if err == nil {
		return errors.New("folder lock already initialized")
	}

	err = os.Mkdir(fileName, os.ModePerm)
	if err != nil {
		return err
	}

	encryptedPassword := m.crypt.Encrypt([]byte(password))

	err = os.WriteFile(filepath.Join(fileName, encryptedDataFileName), encryptedPassword, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (m *windows) Lock() error {
	_, err := os.Stat(fileName)
	if err != nil {
		if err == os.ErrNotExist {
			return errors.New("folder lock not initialized or is still locked")
		}

		return err
	}

	cmd := exec.Command("attrib", "+h", "+s", "+r", fileName)
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (m *windows) Unlock(password string) error {
	_, err := os.Stat(fileName)
	if err != nil {
		if err == os.ErrNotExist {
			return errors.New("folder lock not initialized or is still unlocked")
		}

		return err
	}

	data, err := os.ReadFile(filepath.Join(fileName, encryptedDataFileName))
	if err != nil {
		return err
	}

	decryptedData, err := m.crypt.Decrypt(string(data))
	if err != nil {
		return err
	}

	if password != string(decryptedData) {
		return errors.New("unauthorized")
	}

	cmd := exec.Command("attrib", "-h", "-s", "-r", fileName)
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (m *windows) UpdatePassword(oldPassword, newPassword string) error {
	data, err := os.ReadFile(filepath.Join(fileName, encryptedDataFileName))
	if err != nil {
		return err
	}

	decryptedData, err := m.crypt.Decrypt(string(data))
	if err != nil {
		return err
	}

	if oldPassword != string(decryptedData) {
		return errors.New("unauthorized")
	}

	encryptedPassword := m.crypt.Encrypt([]byte(newPassword))

	err = os.WriteFile(filepath.Join(fileName, encryptedDataFileName), encryptedPassword, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
