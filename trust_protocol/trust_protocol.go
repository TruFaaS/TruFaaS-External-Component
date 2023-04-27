package trust_protocol

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/TruFaaS/TruFaaS/constants"
	"hash"
	"math/big"
	"net/http"
)

type TrustProtocol struct {
	ServerPrivateKey ecdsa.PrivateKey
	ServerPublicKey  ecdsa.PublicKey
	ClientPublicKey  ecdsa.PublicKey
	SharedSecret     big.Int
	curve            elliptic.Curve
	MAC              hash.Hash
}

func (tp *TrustProtocol) GetProtocolInstance(clientPubKeyInBytes []byte) *TrustProtocol {
	// set the default curve
	tp.curve = elliptic.P256()
	// set client public key
	tp.setClientPublicKey(clientPubKeyInBytes)
	//generate server keys
	tp.generateServerKeys()
	// generate secret key
	tp.generateSharedSecret()

	return tp
}

func (tp *TrustProtocol) setClientPublicKey(clientPubKeyInBytes []byte) {
	clientPubKey := &ecdsa.PublicKey{
		Curve: tp.curve,
		X:     new(big.Int).SetBytes(clientPubKeyInBytes[:32]),
		Y:     new(big.Int).SetBytes(clientPubKeyInBytes[32:]),
	}
	tp.ClientPublicKey = *clientPubKey
}

func (tp *TrustProtocol) generateServerKeys() {
	// Generate a private key
	serverPrivKey, _ := ecdsa.GenerateKey(tp.curve, rand.Reader)

	// Get the public key from the private key
	serverPubKey := serverPrivKey.PublicKey

	tp.ServerPrivateKey = *serverPrivKey
	tp.ServerPublicKey = serverPubKey
}

func (tp *TrustProtocol) generateSharedSecret() {
	clientPubKey := tp.ClientPublicKey
	serverPrivKey := tp.ServerPrivateKey
	sharedSecret, _ := clientPubKey.Curve.ScalarMult(clientPubKey.X, clientPubKey.Y, serverPrivKey.D.Bytes())
	tp.SharedSecret = *sharedSecret

}

func (tp *TrustProtocol) GenerateMAC(trustValue string) {

	hMac := hmac.New(sha256.New, tp.SharedSecret.Bytes())
	hMac.Write([]byte(trustValue))
	tp.MAC = hMac
}

func (tp *TrustProtocol) SetResponseHeaders(w http.ResponseWriter, trustValue string) http.ResponseWriter {

	// add trust ca
	w.Header().Set(constants.TrustVerificationHeader, trustValue)

	// Add the MAC tag to the response headers
	macTag := hex.EncodeToString(tp.MAC.Sum(nil))
	w.Header().Set(constants.MACHeader, macTag)

	// Add server's public key to the response headers
	serverPubKeyBytes := append(tp.ServerPublicKey.X.Bytes(), tp.ServerPublicKey.Y.Bytes()...)
	serverPubKeyHex := hex.EncodeToString(serverPubKeyBytes)
	w.Header().Set(constants.ExternalComponentPublicKeyHeader, serverPubKeyHex)

	return w
}
