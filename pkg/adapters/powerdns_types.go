package adapters

import (
	"github.com/joeig/go-powerdns/v3"
	miekgdns "github.com/miekg/dns"
)

func init() {
	powerdnsRRTypeMap = map[uint16]powerdns.RRType{
		miekgdns.TypeA:          powerdns.RRTypeA,
		miekgdns.TypeAAAA:       powerdns.RRTypeAAAA,
		miekgdns.TypeCAA:        powerdns.RRTypeCAA,
		miekgdns.TypeCDNSKEY:    powerdns.RRTypeCDNSKEY,
		miekgdns.TypeCDS:        powerdns.RRTypeCDS,
		miekgdns.TypeCERT:       powerdns.RRTypeCERT,
		miekgdns.TypeCNAME:      powerdns.RRTypeCNAME,
		miekgdns.TypeDHCID:      powerdns.RRTypeDHCID,
		miekgdns.TypeDLV:        powerdns.RRTypeDLV,
		miekgdns.TypeDNAME:      powerdns.RRTypeDNAME,
		miekgdns.TypeDNSKEY:     powerdns.RRTypeDNSKEY,
		miekgdns.TypeDS:         powerdns.RRTypeDS,
		miekgdns.TypeEUI48:      powerdns.RRTypeEUI48,
		miekgdns.TypeEUI64:      powerdns.RRTypeEUI64,
		miekgdns.TypeHINFO:      powerdns.RRTypeHINFO,
		miekgdns.TypeIPSECKEY:   powerdns.RRTypeIPSECKEY,
		miekgdns.TypeKEY:        powerdns.RRTypeKEY,
		miekgdns.TypeKX:         powerdns.RRTypeKX,
		miekgdns.TypeLOC:        powerdns.RRTypeLOC,
		miekgdns.TypeMX:         powerdns.RRTypeMX,
		miekgdns.TypeNAPTR:      powerdns.RRTypeNAPTR,
		miekgdns.TypeNS:         powerdns.RRTypeNS,
		miekgdns.TypeNSEC3:      powerdns.RRTypeNSEC3,
		miekgdns.TypeNSEC3PARAM: powerdns.RRTypeNSEC3PARAM,
		miekgdns.TypeNSEC:       powerdns.RRTypeNSEC,
		miekgdns.TypeOPENPGPKEY: powerdns.RRTypeOPENPGPKEY,
		miekgdns.TypePTR:        powerdns.RRTypePTR,
		miekgdns.TypeRP:         powerdns.RRTypeRP,
		miekgdns.TypeRRSIG:      powerdns.RRTypeRRSIG,
		miekgdns.TypeSIG:        powerdns.RRTypeSIG,
		miekgdns.TypeSMIMEA:     powerdns.RRTypeSMIMEA,
		miekgdns.TypeSOA:        powerdns.RRTypeSOA,
		miekgdns.TypeSPF:        powerdns.RRTypeSPF,
		miekgdns.TypeSRV:        powerdns.RRTypeSRV,
		miekgdns.TypeSSHFP:      powerdns.RRTypeSSHFP,
		miekgdns.TypeTKEY:       powerdns.RRTypeTKEY,
		miekgdns.TypeTLSA:       powerdns.RRTypeTLSA,
		miekgdns.TypeTSIG:       powerdns.RRTypeTSIG,
		miekgdns.TypeTXT:        powerdns.RRTypeTXT,
		miekgdns.TypeURI:        powerdns.RRTypeURI,
	}
}
