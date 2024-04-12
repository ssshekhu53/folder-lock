package mac

import (
	"errors"
	"os"
	"path/filepath"

	"ssshekhu53/folder-lock/services"
)

const (
	fileName              = `private`
	hiddenFileName        = `.private`
	encryptedDataFileName = `.encrypted-data`
)

type mac struct {
	crypt services.Crypt
}

func New(crypt services.Crypt) services.FolderLock {
	return &mac{crypt: crypt}
}

func (m *mac) Init(password string) error {
	_, err := os.Stat(fileName)
	if err == nil {
		return errors.New("folder lock already initialized")
	}

	err = os.Mkdir(fileName, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(fileName, ".nomedia"))
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

func (m *mac) Lock() error {
	_, err := os.Stat(fileName)
	if err != nil {
		if err == os.ErrNotExist {
			return errors.New("folder lock not initialized or is still locked")
		}

		return err
	}

	return os.Rename(fileName, hiddenFileName)
}

func (m *mac) Unlock(password string) error {
	_, err := os.Stat(hiddenFileName)
	if err != nil {
		if err == os.ErrNotExist {
			return errors.New("folder lock not initialized or is still unlocked")
		}

		return err
	}

	data, err := os.ReadFile(filepath.Join(hiddenFileName, encryptedDataFileName))
	if err != nil {
		return err
	}

	decrptedData, err := m.crypt.Decrypt(string(data))
	if err != nil {
		return err
	}

	if password != string(decrptedData) {
		return errors.New("unauthorized")
	}

	return os.Rename(hiddenFileName, fileName)
}

func (m *mac) UpdatePassword(oldPassword, newPassword string) error {
	data, err := os.ReadFile(filepath.Join(hiddenFileName, encryptedDataFileName))
	if err != nil {
		return err
	}

	decrptedData, err := m.crypt.Decrypt(string(data))
	if err != nil {
		return err
	}

	if oldPassword != string(decrptedData) {
		return errors.New("unauthorized")
	}

	encryptedPassword := m.crypt.Encrypt([]byte(newPassword))

	err = os.WriteFile(filepath.Join(fileName, encryptedDataFileName), encryptedPassword, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (m *mac) getFileNameAndAbsPath(path string) (string, string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	}

	fileName := filepath.Base(path)

	return absPath, fileName, nil
}
