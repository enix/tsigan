package adapters

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"github.com/joeig/go-powerdns/v3"
)

const PowerDNSAdapterSlug AdapterSlug = "powerdns"

var powerdnsRRTypeMap map[uint16]powerdns.RRType

func init() {
	registerAdapter(
		PowerDNSAdapterSlug,
		reflect.TypeFor[PowerDNSAdapterConfiguration](),
		reflect.TypeFor[PowerDNSAdapter](),
		NewPowerDNSAdapter)
}

type PowerDNSAdapterConfiguration struct {
	Url        string `validate:"required,http_url"`
	VHost      string `validate:"hostname"`
	Key        string `validate:"base64"`
	decodedKey string
}

type PowerDNSAdapter struct {
	config *PowerDNSAdapterConfiguration
}

func NewPowerDNSAdapter(config IAdapterConfiguration) (IAdapter, error) {
	switch pdnsConfig := config.(type) {
	case *PowerDNSAdapterConfiguration:
		key, error := base64.RawStdEncoding.DecodeString(pdnsConfig.Key)
		if error != nil {
			return nil, fmt.Errorf("failed to decode PowerDNS API key with base64: %w", error)
		}
		pdnsConfig.decodedKey = string(key)
	default:
		panic("PowerDNS adapter created with an unknown config type")
	}

	a := PowerDNSAdapter{config.(*PowerDNSAdapterConfiguration)}
	return &a, nil
}
