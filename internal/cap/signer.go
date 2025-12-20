package cap

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

// Signer handles CAP message digital signing
type Signer struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	algorithm  string
}

// NewSigner creates a new signer from private key file
func NewSigner(privateKeyPath string) (*Signer, error) {
	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}
	
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not RSA key")
		}
	}
	
	return &Signer{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		algorithm:  "RSA-SHA256",
	}, nil
}

// NewSignerFromKey creates a new signer from private key
func NewSignerFromKey(privateKey *rsa.PrivateKey) *Signer {
	return &Signer{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		algorithm:  "RSA-SHA256",
	}
}

// Sign signs a CAP message
func (s *Signer) Sign(capMsg *CAPMessage) error {
	// Convert message to canonical form for signing
	canonical, err := s.canonicalize(capMsg)
	if err != nil {
		return fmt.Errorf("failed to canonicalize message: %w", err)
	}
	
	// Hash the canonical form
	hash := sha256.Sum256([]byte(canonical))
	
	// Sign the hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}
	
	// Encode signature as base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	
	// Add signature to message
	capMsg.Signature = &Signature{
		Algorithm: s.algorithm,
		Value:     signatureBase64,
	}
	
	return nil
}

// Verify verifies a CAP message signature
func (s *Signer) Verify(capMsg *CAPMessage) error {
	if capMsg.Signature == nil {
		return fmt.Errorf("message has no signature")
	}
	
	// Remove signature for verification
	sigCopy := capMsg.Signature
	capMsg.Signature = nil
	
	// Convert message to canonical form
	canonical, err := s.canonicalize(capMsg)
	if err != nil {
		return fmt.Errorf("failed to canonicalize message: %w", err)
	}
	
	// Hash the canonical form
	hash := sha256.Sum256([]byte(canonical))
	
	// Decode signature
	signature, err := base64.StdEncoding.DecodeString(sigCopy.Value)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}
	
	// Verify signature
	err = rsa.VerifyPKCS1v15(s.publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}
	
	// Restore signature
	capMsg.Signature = sigCopy
	
	return nil
}

// canonicalize converts CAP message to canonical form for signing
func (s *Signer) canonicalize(capMsg *CAPMessage) (string, error) {
	// Simple canonicalization: use XML representation (without signature)
	xml, err := capMsg.ToXML()
	if err != nil {
		return "", err
	}
	
	// Remove signature element if present (for canonicalization)
	// In production, would use proper XML canonicalization (C14N)
	return xml, nil
}

