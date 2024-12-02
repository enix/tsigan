package update

import (
	"fmt"

	miekgdns "github.com/miekg/dns"
	"go.uber.org/zap"
)

func Execute(authorization *Authorization, prerequisites *Prerequisites, update *[]miekgdns.RR,
	logger *zap.SugaredLogger) error {

	// Validate authorizations
	if err := authorization.Evaluate(); err != nil {
		return fmt.Errorf("authorization failed: %w", err)
	}

	zone := authorization.Zone()
	adapter := zone.GetHandler()

	// Start an adapter transaction
	transaction, err := adapter.NewTransaction()
	if err != nil {
		return fmt.Errorf("new adapter transaction: %w", err)
	}

	// Validate all update prerequisites
	if err := prerequisites.Evaluate(transaction); err != nil {
		return fmt.Errorf("prerequisites failed: %w", err)
	}

	if err := transaction.Commit(update); err != nil {
		return fmt.Errorf("transaction failure: %w", err)
	}

	return nil
}
