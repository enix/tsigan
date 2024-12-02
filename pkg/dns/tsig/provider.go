package tsig

import (
	"crypto/hmac"
	"encoding/hex"

	miekgdns "github.com/miekg/dns"
	"go.uber.org/zap"
)

type TsigProvider struct {
	keyring *TsigKeyring
	logger  *zap.SugaredLogger
}

func NewTsigProvider(keyring *TsigKeyring, logger *zap.SugaredLogger) *TsigProvider {
	return &TsigProvider{keyring, logger}
}

func (p *TsigProvider) generate(msg []byte, t *miekgdns.TSIG) ([]byte, error) {
	keyName := t.Hdr.Name

	key := p.keyring.GetKey(keyName)
	if key == nil {
		p.logger.Debugw("signature generation failure: key name is unknown", "keyname", keyName)
		return nil, miekgdns.ErrSecret
	}

	// TODO check if canonicalization of t.Algorithm is needed
	tsigHmac, err := NewHmac(t.Algorithm)
	if err != nil {
		p.logger.Debugw("signature generation failure: invalid algorithm",
			"keyname", keyName, "algorithm", t.Algorithm, "error", err.Error())
		return nil, miekgdns.ErrKeyAlg
	}

	return tsigHmac.Sum(msg, key)
}

func (p *TsigProvider) Generate(msg []byte, t *miekgdns.TSIG) ([]byte, error) {
	keyName := t.Hdr.Name
	p.logger.Debugw("message signature generation request", "keyname", keyName)
	return p.generate(msg, t)
}

func (p *TsigProvider) Verify(msg []byte, t *miekgdns.TSIG) error {
	keyName := t.Hdr.Name
	p.logger.Debugw("message signature verification request", "keyname", keyName)

	computedMac, err := p.generate(msg, t)
	if err != nil {
		p.logger.Debugw("signature verification failure: cannot generate expected MAC",
			"keyname", keyName, "error", err.Error())
		return err
	}

	receivedMac, err := hex.DecodeString(t.MAC)
	if err != nil {
		p.logger.Debugw("signature verification failure: cannot decode received MAC",
			"keyname", keyName, "error", err.Error())
		return err
	}

	if !hmac.Equal(computedMac, receivedMac) {
		p.logger.Debugw("message signature verification failed, MACs are not equal", "keyname", keyName)
		return miekgdns.ErrSig
	}
	return nil
}
