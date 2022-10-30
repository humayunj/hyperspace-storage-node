package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/fatih/color"
)

type StoredFile struct {
	Hash         []byte
	Path         string
	AddedAt      uint64
	OwnerAddress string
}

type FileService struct {
	db *badger.DB
}

func (f *FileService) Open() {
	db, err := badger.Open(badger.DefaultOptions("files.db"))
	if err != nil {
		log.Fatal("Failed to open db")
	}

	f.db = db
}

func (f *FileService) Close() {
	f.db.Close()
	f.db = nil
}

func (f *FileService) AddFile(r io.Reader, key []byte, ownerAddress string) error {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		log.Fatal(err)
		return err
	}
	hash := h.Sum(nil)
	path := strings.Join([]string{"uploads/", hex.EncodeToString(hash)}, "")
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, r); err != nil {
		log.Fatal(err)
		return err
	}

	var sf StoredFile
	sf.Hash = hash
	sf.OwnerAddress = ownerAddress
	sf.AddedAt = uint64(time.Now().Unix())
	sf.Path = path

	sfBytes, err := sf.Encode()
	if err != nil {
		log.Fatal(err)
		os.Remove(path)
		return err
	}

	err = f.db.Update(func(txn *badger.Txn) error {
		err = txn.Set(bytes.Join([][]byte{[]byte("file_"), hash}, []byte{}), sfBytes)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
		os.Remove(path)
	} else {
		color.Green("Added File: " + hex.EncodeToString(hash))
	}
	return err
}

func (f *FileService) GetAllFiles() ([]*StoredFile, error) {
	files := make([]*StoredFile, 0, 10)

	err := f.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.ValueCopy(nil)
			if err != nil {
				log.Fatal(err)
				continue
			}
			if !bytes.Equal(k[:5], []byte("file_")) {
				continue
			}
			sf := StoredFile{}
			sf.Decode(v)
			files = append(files, &sf)
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (f *StoredFile) Decode(b []byte) error {
	buff := bytes.Buffer{}
	buff.Write(b)
	decoder := gob.NewDecoder(&buff)
	return decoder.Decode(f)

}

func (f *StoredFile) Encode() ([]byte, error) {
	buff := bytes.Buffer{}

	encoder := gob.NewEncoder(&buff)
	if err := encoder.Encode(f); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
