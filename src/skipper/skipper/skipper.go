package skipper

import "net/http"

type RawData interface {

	// todo: consider what config format failures should be considered invalid
	// in json:
	// {
	//     "backends": {"pdp": "https://www.zalando.de/pdp.html"},
	//     "frontends": [{
	//         "route": "PathRegexp(`.*\\.html`)",
	//         "backendId": "pdp",
	//         "filters": [
	//             "pdp-custom-headers",
	//             "x-session-id"
	//         ]
	//     }],
	//     "filter-specs": {
	//         "pdp-custom-headers": {
	//             "middleware-name": "custom-headers",
	//             "config": {"free-data": 3.14}
	//         },
	//         "x-session-id": {
	//             "middleware-name": "x-session-id",
	//             "config": {"generator": "v4"}
	//         }
	//     }
	// }
	Get() map[string]interface{}
}

type DataClient interface {
	Receive() <-chan RawData
}

type Backend interface {
	Url() string
}

type Filter interface {
	http.Handler
	Id() string
}

type Route interface {
	Backend() Backend
	Filters() []Filter
}

type Settings interface {
	Route(*http.Request) (Route, error)
	Address() string
}

type SettingsSource interface {
	Subscribe(chan<- Settings)
}

type SettingsDispatcher interface {
	SettingsSource
	Push() chan<- Settings
}

type SettingsProcessor interface {
	SettingsSource
	PushRawData() chan<- RawData
}

type MiddlewareConfig map[string]interface{}

type Middleware interface {
	Name() string
	MakeFilter(id string, s MiddlewareConfig) Filter
}

type MiddlewareRegistry interface {
	Add(...Middleware)
	Get(string) Middleware
	Remove(string)
}
