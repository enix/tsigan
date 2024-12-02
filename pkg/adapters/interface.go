package adapters

type IAdapterConfiguration interface{}

type IAdapter interface {
	// NewTransaction() *IAdapterTransaction
}

type IAdapterTransaction interface {
	// Execute() error
}
