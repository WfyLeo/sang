package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"

	"github.com/forgoer/openssl"
)

// MD5Hash returns the MD5 hash of the given data using the openssl package.
func MD5Hash(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// Encrypt encrypts the given data using AES encryption.
func Encrypt(data, key, iv string) (string, error) {
	encrypted, err := openssl.AesCBCEncrypt([]byte(data), []byte(key), []byte(iv), openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(encrypted), nil
}

// Decrypt decrypts the given AES encrypted data.
func Decrypt(encryptedData, key, iv string) (string, error) {
	data, err := hex.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	decrypted, err := openssl.AesCBCDecrypt(data, []byte(key), []byte(iv), openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// Compress compresses the given data using zlib.
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Decompress decompresses the given zlib compressed data.
func Decompress(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	result, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// WriteUint32 writes a uint32 value to a byte slice using little-endian encoding.
func WriteUint32(value uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ReadUint32 reads a uint32 value from a byte slice using little-endian encoding.
func ReadUint32(data []byte) (uint32, error) {
	buf := bytes.NewReader(data)
	var value uint32
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}
