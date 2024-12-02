package server

import (
	"crypto/hmac"
	"encoding/hex"

	"github.com/enix/tsigan/pkg/tsig"
	miekgdns "github.com/miekg/dns"
)

type TsigProvider struct {
	keyring *tsig.TsigKeyring
}

func NewTsigProvider(keyring *tsig.TsigKeyring) *TsigProvider {
	return &TsigProvider{keyring}
}

func (p *TsigProvider) Generate(msg []byte, t *miekgdns.TSIG) ([]byte, error) {
	keyName := t.Hdr.Name
	Logger.Debugw("message signature request", "keyname", keyName)
	return p.generate(msg, t)
}

func (p *TsigProvider) generate(msg []byte, t *miekgdns.TSIG) ([]byte, error) {
	keyName := t.Hdr.Name

	key := p.keyring.GetKey(keyName)
	if key == nil {
		Logger.Debugw("cannot sign, key name is unknown", "keyname", keyName)
		return nil, miekgdns.ErrSecret
	}

	// TODO check if canonicalization of t.Algorithm is needed
	tsigHmac, err := tsig.NewHmac(t.Algorithm)
	if err != nil {
		Logger.Debugw("cannot sign, invalid algorithm",
			"keyname", keyName, "algorithm", t.Algorithm, "error", err.Error())
		return nil, miekgdns.ErrKeyAlg
	}

	return tsigHmac.Sum(msg, key)
}

func (p *TsigProvider) Verify(msg []byte, t *miekgdns.TSIG) error {
	keyName := t.Hdr.Name
	Logger.Debugw("message signature verification request", "keyname", keyName)

	computedMac, err := p.generate(msg, t)
	if err != nil {
		Logger.Debugw("cannot check signature, failed to generate expectations",
			"keyname", keyName, "error", err.Error())
		return err
	}

	receivedMac, err := hex.DecodeString(t.MAC)
	if err != nil {
		Logger.Debugw("cannot check signature, failed to decode received MAC",
			"keyname", keyName, "error", err.Error())
		return err
	}

	if !hmac.Equal(computedMac, receivedMac) {
		Logger.Debugw("message signature verification failed, MACs not equal", "keyname", keyName)
		return miekgdns.ErrSig
	}
	return nil
}
