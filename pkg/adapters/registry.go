package adapters

import (
	"fmt"
	"reflect"
)

var adapters = map[AdapterSlug]*adapterInfo{}

type AdapterSlug string

type adapterFactory func(IAdapterConfiguration) (IAdapter, error)

type adapterInfo struct {
	slug         AdapterSlug
	concreteType reflect.Type
	configType   reflect.Type
	factory      adapterFactory
}

func (s AdapterSlug) IsValid() bool {
	_, found := adapters[s]
	return found
}

func registerAdapter(slug AdapterSlug, concreteType reflect.Type, configType reflect.Type, factory adapterFactory) {
	if _, found := adapters[slug]; found {
		panic("adapter slug collision")
	}
	adapters[slug] = &adapterInfo{slug, concreteType, configType, factory}
}

func getAdapterInfo(slug AdapterSlug) (*adapterInfo, error) {
	info, found := adapters[slug]
	if !found {
		return nil, fmt.Errorf("invalid adapter type '%s'", slug)
	}
	return info, nil
}

func NewAdapterConfiguration(slug AdapterSlug) (IAdapterConfiguration, error) {
	info, err := getAdapterInfo(slug)
	if err != nil {
		return nil, err
	}

	newConfig := reflect.New(info.configType).Interface()
	return newConfig, nil
}

func NewAdapter(slug AdapterSlug, configuration IAdapterConfiguration) (IAdapter, error) {
	var (
		err     error
		adapter IAdapter
	)

	info, err := getAdapterInfo(slug)
	if err != nil {
		return nil, err
	}

	adapter, err = info.factory(configuration)
	if err != nil {
		panic("FIXME message")
	}

	fmt.Printf("[NewAdapter] adapter=%p\n", adapter)
	return adapter, nil
}
