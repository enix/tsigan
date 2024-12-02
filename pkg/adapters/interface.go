package adapters

import (
	miekgdns "github.com/miekg/dns"
)

type IAdapterConfiguration interface{}

type IAdapter interface {
	NewTransaction() (IAdapterTransaction, error)
}

type IAdapterTransaction interface {
	Commit(*[]miekgdns.RR) error
}
