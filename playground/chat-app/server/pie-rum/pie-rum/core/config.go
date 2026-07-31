package pierum

import (
	"sync"
)

type IConfig struct {
	Activate     bool     `json:"active"`
	SwapOverview *ISwitch `json:"switch"`
	mu           sync.Mutex
}

// defaultConfig returns the active config by default
func defaultConfig() *IConfig {
	return &IConfig{
		Activate: true,
		SwapOverview: &ISwitch{
			HSwitch: false,
			Name:    "",
		},
	}
}

type ISwitch struct {
	Name    string `json:"switched_with"` // name that is switch with
	HSwitch bool   `json:"is_switch"`     // has_switch
}

func (i *IConfig) setSwapOverview(info *ISwitch) {
	i.SwapOverview = info
}

func (i *IConfig) GetSwapOverview() *ISwitch {
	return i.SwapOverview
}

func (i *IConfig) setActivate(status bool) {
	i.Activate = status
}

func (i *IConfig) getActivate() bool {
	return i.Activate
}

// IConfigRequest implemets the configuration changes required for profiles, services, and events
// profile is mandatory
// usage: profile only-> changes will be made for profile
// profile->service changes will be made for service
// profile->service->event changes will be made for event
// type IConfigRequest struct {
// 	Profile, Kit, Service, Dispatcher, Event string

//		Swap string // to swap with
//	}

type IConfigRequest struct {
	Action                                   string // "activate", "deactivate", "swap"
	Target                                   string // "profile", "kit", "service", "dispatcher", "event"
	Profile, Kit, Service, Dispatcher, Event string
	Config                                   *IConfig // set the config directly
	Swap                                     string
}
