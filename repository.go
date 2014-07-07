package main

import (
	"bufio"
	"io"
	"os"
	"path"
)

const (
	idFilename     = ".gpg-id"
	fileExtension  = ".gpg"
	filePermission = 0600
	dirPermission  = 0700
)

type Repository interface {
	Get(key string, passphrase []byte) (io.ReadCloser, error)
	Put(key string) (io.WriteCloser, error)
}

type fileRepository struct {
	root string
}

func NewRepository(root string) Repository {
	return &fileRepository{
		root: root,
	}
}

func (r fileRepository) Ids() ([]string, error) {
	file, err := os.Open(path.Join(r.root, idFilename))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ids := []string{}
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	return ids, scanner.Err()
}

func (r fileRepository) Get(key string, passphrase []byte) (io.ReadCloser, error) {
	ids, err := r.Ids()
	if err != nil {
		return nil, err
	}
	el, err := EntitiesFromKeyRing(SecureKeyRingName(), ids)
	if err != nil {
		return nil, err
	}
	return OpenEncrypted(path.Join(r.root, key+fileExtension), el, passphrase)
}

func (r fileRepository) Put(key string) (io.WriteCloser, error) {
	ids, err := r.Ids()
	if err != nil {
		return nil, err
	}
	el, err := EntitiesFromKeyRing(PublicKeyRingName(), ids)
	if err != nil {
		return nil, err
	}
	filepath := path.Join(r.root, key+fileExtension)
	if err := os.MkdirAll(path.Dir(filepath), dirPermission); err != nil {
		return nil, err
	}
	return CreateEncrypted(filepath, filePermission, el)
}
