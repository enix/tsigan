package server

import (
	"fmt"

	"github.com/enix/tsigan/pkg/adapters"
)

func (s *Server) init() error {
	Logger.Debug("initializing server state")

	// process TSIG keys from configuration
	Logger.Debugw("initializing server keyring", "count", len(s.Configuration.Tsig.Keys))
	for _, config := range s.Configuration.Tsig.Keys {
		if err := s.newKey(&config); err != nil {
			return err
		}
	}
	Logger.Debug("finished initializing server keyring")

	// process handlers from configuration
	Logger.Debugw("initializing server handler adapters", "count", len(s.Configuration.Handlers))
	for _, config := range s.Configuration.Handlers {
		if err := s.newHandler(&config); err != nil {
			return err
		}
	}
	Logger.Debug("finished initializing server handler adapters")

	// process zones from configuration
	Logger.Debugw("initializing server zones", "count", len(s.Configuration.Zones))
	for _, config := range s.Configuration.Zones {
		if err := s.newZone(&config); err != nil {
			return err
		}
	}
	Logger.Debug("finished initializing server zones")

	Logger.Debug("finished initializing server state")
	return nil
}

func (s *Server) newKey(config *TsigKeyConfiguration) error {
	Logger.Debugw("initializing new key", "name", config.Name)

	if err := s.keyring.AddEncodedKey(config.Name, config.Key); err != nil {
		return fmt.Errorf("adding key '%s' to keyring: %w", config.Name, err)
	}

	if config.Default {
		if len(s.defaultKeyName) > 0 {
			// to force code refactoring when config reload is implemented
			// this would otherwise cause auth bugs with zones using a default key
			Logger.Fatal("unsupported default key reset attempted")
		}
		Logger.Debugw("promoting key as default", "name", config.Name)
		s.defaultKeyName = config.Name
	}
	return nil
}

func (s *Server) newHandler(config *HandlerConfiguration) error {
	Logger.Debugw("initializing new handler adapter", "name", config.Name)

	adapter, err := adapters.NewAdapter(config.Adapter, config.Settings)
	if err != nil {
		return fmt.Errorf("adapter '%s': %w", config.Name, err)
	}

	s.adapters = append(s.adapters, &adapter)
	s.adaptersByName[config.Name] = &adapter

	Logger.Debugw("initialized new handler adapter", "name", config.Name, "object", fmt.Sprintf("%p", &adapter))

	if config.Default {
		if s.defaultAdapter != nil {
			// to force code refactoring when config reload is implemented
			// this would otherwise cause bugs with zones using a default handler
			Logger.Fatal("can't have more than one default handler")
		}
		s.defaultAdapter = &adapter
		Logger.Debugw("handler adapter promoted to default", "name", config.Name, "object", fmt.Sprintf("%p", &adapter))
	}

	return nil
}

func (s *Server) newZone(config *ZoneConfiguration) error {
	Logger.Debugw("initializing new zone", "name", config.Zone)

	zone, err := NewZone(config.Zone)
	if err != nil {
		Logger.Fatalw("failed to initialize new zone", "name", config.Zone, "error", err.Error())
	}

	if config.Unsecure == false {
		// auth enabled, processing keys
		Logger.Debugw("zone has authentication enabled", "zone", config.Zone)

		addKeys := make([]string, 0)
		if len(config.Keys) > 0 {
			// add keys from zone config
			addKeys = append(addKeys, config.Keys...)
		} else {
			// or try adding the default key
			if len(s.defaultKeyName) > 0 {
				addKeys = append(addKeys, s.defaultKeyName)
			}
		}

		if len(addKeys) == 0 {
			Logger.Fatalw("zone with authentication enabled but no key", "zone", config.Zone)
		}

		// push keys to zone
		for _, key := range addKeys {
			if s.keyring.HasKey(key) {
				zone.AddValidKey(key)
			} else {
				Logger.Fatalw("zone requesting unknown key", "zone", config.Zone, "keyname", key)
			}
		}
	} else {
		Logger.Warnw("zone has authentication disabled", "zone", config.Zone)
		zone.DisableAuthentication()
	}

	s.zones = append(s.zones, zone)
	s.zonesByFqdn[zone.GetFqdn()] = zone
	return nil
}
