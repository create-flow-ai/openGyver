package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"strings"
	"testing"
)

func TestIsHexKey(t *testing.T) {
	valid := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if !isHexKey(valid) {
		t.Errorf("expected isHexKey(%q) = true", valid)
	}

	if isHexKey("abcdef") {
		t.Error("expected isHexKey for short string = false")
	}

	invalid := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
	if isHexKey(invalid) {
		t.Error("expected isHexKey for non-hex string = false")
	}
}

func TestDeriveKey(t *testing.T) {
	salt := []byte("testsalt12345678")
	key := deriveKey([]byte("my passphrase"), salt)
	if len(key) != 32 {
		t.Errorf("expected 32-byte key, got %d bytes", len(key))
	}

	// Deterministic.
	key2 := deriveKey([]byte("my passphrase"), salt)
	if hex.EncodeToString(key) != hex.EncodeToString(key2) {
		t.Error("deriveKey is not deterministic")
	}

	// Different passphrase => different key.
	key3 := deriveKey([]byte("other passphrase"), salt)
	if hex.EncodeToString(key) == hex.EncodeToString(key3) {
		t.Error("different passphrases produced same key")
	}
}

func TestAESEncryptDecryptRoundtrip_Passphrase(t *testing.T) {
	passphrase := "test-passphrase"
	plaintext := "hello, world! this is a secret message"

	// Encrypt.
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("generating salt: %v", err)
	}
	keyBytes := deriveKey([]byte(passphrase), salt)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		t.Fatalf("creating cipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("creating GCM: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		t.Fatalf("generating nonce: %v", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	blob := append(salt, ciphertext...)
	encoded := base64.StdEncoding.EncodeToString(blob)

	// Decrypt.
	decodedBlob, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}
	decodedSalt := decodedBlob[:16]
	decodedCiphertext := decodedBlob[16:]
	decKeyBytes := deriveKey([]byte(passphrase), decodedSalt)

	block2, err := aes.NewCipher(decKeyBytes)
	if err != nil {
		t.Fatalf("creating cipher for decrypt: %v", err)
	}
	gcm2, err := cipher.NewGCM(block2)
	if err != nil {
		t.Fatalf("creating GCM for decrypt: %v", err)
	}
	nonceSize := gcm2.NonceSize()
	decNonce, decCt := decodedCiphertext[:nonceSize], decodedCiphertext[nonceSize:]
	decrypted, err := gcm2.Open(nil, decNonce, decCt, nil)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if string(decrypted) != plaintext {
		t.Errorf("roundtrip mismatch: got %q, want %q", string(decrypted), plaintext)
	}
}

func TestAESEncryptDecryptRoundtrip_HexKey(t *testing.T) {
	hexKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	plaintext := "secret data with hex key"

	keyBytes, _ := hex.DecodeString(hexKey)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		t.Fatalf("creating cipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("creating GCM: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		t.Fatalf("generating nonce: %v", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	// Decrypt.
	blob, _ := base64.StdEncoding.DecodeString(encoded)
	block2, _ := aes.NewCipher(keyBytes)
	gcm2, _ := cipher.NewGCM(block2)
	nonceSize := gcm2.NonceSize()
	decNonce, decCt := blob[:nonceSize], blob[nonceSize:]
	decrypted, err := gcm2.Open(nil, decNonce, decCt, nil)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}
	if string(decrypted) != plaintext {
		t.Errorf("roundtrip mismatch: got %q, want %q", string(decrypted), plaintext)
	}
}

func TestAESDecryptWrongKey(t *testing.T) {
	passphrase := "correct-passphrase"
	plaintext := "secret"

	salt := make([]byte, 16)
	rand.Read(salt)
	keyBytes := deriveKey([]byte(passphrase), salt)
	block, _ := aes.NewCipher(keyBytes)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	blob := append(salt, ciphertext...)

	// Try to decrypt with wrong passphrase.
	wrongKey := deriveKey([]byte("wrong-passphrase"), salt)
	block2, _ := aes.NewCipher(wrongKey)
	gcm2, _ := cipher.NewGCM(block2)
	nonceSize := gcm2.NonceSize()
	ct := blob[16:]
	decNonce, decCt := ct[:nonceSize], ct[nonceSize:]
	_, err := gcm2.Open(nil, decNonce, decCt, nil)
	if err == nil {
		t.Error("expected decryption to fail with wrong key")
	}
}

func TestRSAKeyGeneration(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating RSA key: %v", err)
	}

	privDER, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		t.Fatalf("marshalling private key: %v", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})

	pubDER, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		t.Fatalf("marshalling public key: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})

	if !strings.Contains(string(privPEM), "-----BEGIN PRIVATE KEY-----") {
		t.Error("private PEM missing expected header")
	}
	if !strings.Contains(string(privPEM), "-----END PRIVATE KEY-----") {
		t.Error("private PEM missing expected footer")
	}
	if !strings.Contains(string(pubPEM), "-----BEGIN PUBLIC KEY-----") {
		t.Error("public PEM missing expected header")
	}
	if !strings.Contains(string(pubPEM), "-----END PUBLIC KEY-----") {
		t.Error("public PEM missing expected footer")
	}
}

func TestWriteKeyFiles(t *testing.T) {
	dir := t.TempDir()
	data1 := []byte("file-one-content")
	data2 := []byte("file-two-content")

	err := writeKeyFiles(dir, "a.pem", data1, "b.pem", data2)
	if err != nil {
		t.Fatalf("writeKeyFiles error: %v", err)
	}
}
