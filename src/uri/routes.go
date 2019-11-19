package uri

import (
	"net/http"

	"uri/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"HandleGetBlockchain",
		"GET",
		"/",
		handlers.HandleGetBlockchain,
	},
	Route{
		"HandlePeers",
		"GET",
		"/peer",
		handlers.HandlePeers,
	},
	Route{
		"HandleBlocks",
		"POST",
		"/block/{height}/{hash}",
		handlers.HandleBlocks,
	},
	Route{
		"HandleShow",
		"GET",
		"/show",             // sample => askoddoreven/5
		handlers.HandleShow, //api to send post request
	},
	Route{
		"HandleUpload",
		"GET",
		"/upload",
		handlers.HandleUpload,
	},
	Route{
		"HandleHeartbeatReceive",
		"POST",
		"/heartbeat/receive",
		handlers.HandleHeartbeatReceive,
	},
	Route{
		"HandleStart",
		"GET",
		"/start",
		handlers.HandleStart,
	},
}
