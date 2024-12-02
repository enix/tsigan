package adapters

import (
	"fmt"
	"reflect"
)

var adapters = map[AdapterSlug]*adapterInfo{}

type AdapterSlug string

type adapterFactory func(IAdapterConfiguration) (IAdapter, error)
type transactionFactory func(IAdapter) (IAdapterTransaction, error)

type adapterInfo struct {
	slug         AdapterSlug
	configType   reflect.Type
	concreteType reflect.Type
	factory      adapterFactory
}

func (s AdapterSlug) IsValid() bool {
	_, found := adapters[s]
	return found
}

func registerAdapter(slug AdapterSlug, configType reflect.Type, concreteType reflect.Type,
	factory adapterFactory) {

	if _, found := adapters[slug]; found {
		panic("adapter slug collision")
	}
	adapters[slug] = &adapterInfo{
		slug,
		configType,
		concreteType,
		factory,
	}
}

func adapterInfoBySlug(slug AdapterSlug) (*adapterInfo, error) {
	info, found := adapters[slug]
	if !found {
		return nil, fmt.Errorf("invalid adapter type '%s'", slug)
	}
	return info, nil
}

func adapterInfoByConfigType(configType reflect.Type) (*adapterInfo, error) {
	slug := AdapterSlug("")
	for k, adapter := range adapters {
		if configType == adapter.configType {
			slug = k
			break
		}
	}
	return adapterInfoBySlug(slug)
}

func adapterInfoByAdapterType(concreteType reflect.Type) (*adapterInfo, error) {
	slug := AdapterSlug("")
	for k, adapter := range adapters {
		if concreteType == adapter.concreteType {
			slug = k
			break
		}
	}
	return adapterInfoBySlug(slug)
}

func NewAdapterConfiguration(slug AdapterSlug) (IAdapterConfiguration, error) {
	info, err := adapterInfoBySlug(slug)
	if err != nil {
		return nil, err
	}

	newConfig := reflect.New(info.configType).Interface()
	return newConfig, nil
}

func NewAdapter(configuration IAdapterConfiguration) (IAdapter, error) {
	info, err := adapterInfoByConfigType(reflect.TypeOf(configuration).Elem())
	if err != nil {
		return nil, err
	}
	adapter, err := info.factory(configuration)
	if err != nil {
		panic("FIXME message")
	}

	return adapter, nil
}
