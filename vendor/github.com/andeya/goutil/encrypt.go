package goutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash/fnv"
	"io"

	"github.com/henrylee2cn/ameda"
)

// Md5 returns the MD5 checksum string of the data.
func Md5(b []byte) string {
	checksum := md5.Sum(b)
	return hex.EncodeToString(checksum[:])
}

// Sha1 returns the sha1 checksum string of the data.
func Sha1(b []byte) string {
	checksum := sha1.Sum(b)
	return hex.EncodeToString(checksum[:])
}

// Sha256 returns the sha256 checksum string of the data.
func Sha256(b []byte) string {
	checksum := sha256.Sum256(b)
	return hex.EncodeToString(checksum[:])
}

// Sha512 returns the sha512 checksum string of the data.
func Sha512(b []byte) string {
	checksum := sha512.Sum512(b)
	return hex.EncodeToString(checksum[:])
}

// Fnv1aToUint64 returns the 64-bit FNV-1a hash sum of b.
func Fnv1aToUint64(b []byte) uint64 {
	h := fnv.New64a()
	h.Reset()
	h.Write(b)
	return h.Sum64()
}

// Fnv1aToUint32 returns the 32-bit FNV-1a hash sum of b.
func Fnv1aToUint32(b []byte) uint32 {
	h := fnv.New32a()
	h.Reset()
	h.Write(b)
	return h.Sum32()
}

// AESEncrypt uses ECB mode to encrypt a piece of data and then encodes it in hex.
// The cipherkey argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func AESEncrypt(cipherkey, plaintext []byte, useBase64 ...bool) []byte {
	block := mustNewCipher(cipherkey)
	blockSize := block.BlockSize()
	plaintext = pkcs5Padding(plaintext, blockSize)
	r := make([]byte, len(plaintext))
	dst := r
	for len(plaintext) > 0 {
		block.Encrypt(dst, plaintext)
		plaintext = plaintext[blockSize:]
		dst = dst[blockSize:]
	}
	return encode(r, useBase64)
}

// AESDecrypt hex decodes a piece of data and then decrypts it using ECB mode.
// The cipherkey argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func AESDecrypt(cipherkey, ciphertext []byte, useBase64 ...bool) ([]byte, error) {
	src, err := decode(ciphertext, useBase64)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(cipherkey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	r := make([]byte, len(src))
	dst := r
	for len(src) > 0 {
		block.Decrypt(dst, src)
		src = src[blockSize:]
		dst = dst[blockSize:]
	}
	return pkcs5Unpadding(r)
}

// AESCBCEncrypt uses CBC mode to encrypt a piece of data and then encodes it in hex.
// The cipherkey argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func AESCBCEncrypt(cipherkey, plaintext []byte, useBase64 ...bool) []byte {
	block := mustNewCipher(cipherkey)
	blockSize := block.BlockSize()
	plaintext = pkcs5Padding(plaintext, blockSize)
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return encode(ciphertext, useBase64)
}

// AESCBCDecrypt hex decodes a piece of data and then decrypts it using CBC mode.
// The cipherkey argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func AESCBCDecrypt(cipherkey, ciphertext []byte, useBase64 ...bool) ([]byte, error) {
	ciphertext, err := decode(ciphertext, useBase64)
	if err != nil {
		return nil, err
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	block, err := aes.NewCipher(cipherkey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	plaintext := ciphertext
	mode.CryptBlocks(plaintext, ciphertext)
	return pkcs5Unpadding(plaintext)
}

// AESCTREncrypt uses CTR mode to encrypt a piece of data and then encodes it in hex.
// The cipherkey argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func AESCTREncrypt(cipherkey, plaintext []byte, useBase64 ...bool) []byte {
	block := mustNewCipher(cipherkey)
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return encode(ciphertext, useBase64)
}

// AESCTRDecrypt hex decodes a piece of data and then decrypts it using CTR mode.
// The cipherkey argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func AESCTRDecrypt(cipherkey, ciphertext []byte, useBase64 ...bool) ([]byte, error) {
	ciphertext, err := decode(ciphertext, useBase64)
	if err != nil {
		return nil, err
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	block, err := aes.NewCipher(cipherkey)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	plaintext := ciphertext
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

func mustNewCipher(cipherkey []byte) cipher.Block {
	block, err := aes.NewCipher(cipherkey)
	if err != nil {
		panic(err)
	}
	return block
}

func pkcs5Padding(plaintext []byte, blockSize int) []byte {
	n := byte(blockSize - len(plaintext)%blockSize)
	for i := byte(0); i < n; i++ {
		plaintext = append(plaintext, n)
	}
	return plaintext
}

func pkcs5Unpadding(r []byte) ([]byte, error) {
	l := len(r)
	if l == 0 {
		return nil, errors.New("input padded bytes is empty")
	}
	last := int(r[l-1])
	if l-last < 0 {
		return nil, errors.New("input padded bytes is invalid")
	}
	n := byte(last)
	pad := r[l-last : l]
	isPad := true
	for _, v := range pad {
		if v != n {
			isPad = false
			break
		}
	}
	if !isPad {
		return nil, errors.New("remove pad error")
	}
	return r[:l-last], nil
}

func encode(src []byte, useBase64 []bool) []byte {
	if ameda.OneBool(useBase64) {
		return base64Encode(src)
	}
	return hexEncode(src)
}

func decode(src []byte, useBase64 []bool) ([]byte, error) {
	if ameda.OneBool(useBase64) {
		return base64Decode(src)
	}
	return hexDecode(src)
}

func hexEncode(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

func hexDecode(src []byte) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	return dst, err
}

func base64Encode(src []byte) []byte {
	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(src)))
	base64.RawURLEncoding.Encode(buf, src)
	return buf
}

func base64Decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(src)))
	n, err := base64.RawURLEncoding.Decode(dst, src)
	return dst[:n], err
}
