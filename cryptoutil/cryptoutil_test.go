package cryptoutil

import (
	"testing"	
)

func TestAES128ECB(t *testing.T) {
	source := []byte("abcdefghi 123456")
	key := []byte("1234567890123456")
	enc := AES128ECBEncrypt(source, key)
	dec := AES128ECBDecrypt(enc, key)
	if string(dec) != string(source) {
		t.Errorf("%s is different from %s", source, dec)
	}
}

func TestAES128CBC(t *testing.T) {
	source := []byte("abcdefghi 123456abcdefghi 123456abcdefghi 123456")
	key := []byte("1234567890123456")
	iv := []byte("7777777777777777")
	enc := AES128CBCEncrypt(source, key, iv)
	dec := AES128CBCDecrypt(enc, key, iv)
	if string(dec) != string(source) {
		t.Errorf("%s is different from %s", source, dec)
	}
}