package server

import (
	"fmt"

	"github.com/enix/tsigan/pkg/adapters"
	miekgdns "github.com/miekg/dns"
)

type Zone struct {
	fqdn      string
	handler   *adapters.IAdapter
	validKeys []string
	unsecure  bool
}

func NewZone(name string) (*Zone, error) {
	fqdn := miekgdns.Fqdn(miekgdns.CanonicalName(name))
	if len(fqdn) == 0 {
		return nil, fmt.Errorf("zone FQDN was empty after canonicalization")
	}

	return &Zone{
		fqdn:     fqdn,
		handler:  nil,
		unsecure: false,
	}, nil
}

func (z *Zone) GetFqdn() string {
	return z.fqdn
}

func (z *Zone) SetHandler(adapter *adapters.IAdapter) {
	z.handler = adapter
}

func (z *Zone) AddValidKey(name string) {
	z.validKeys = append(z.validKeys, name)
}

func (z *Zone) KeyIsAuthorized(name string) bool {
	for _, k := range z.validKeys {
		if k == name {
			return true
		}
	}
	return false
}

func (z *Zone) DisableAuthentication() {
	z.unsecure = true
	z.validKeys = nil
}

func (z *Zone) HasAuthenticationDisabled() bool {
	return z.unsecure
}
