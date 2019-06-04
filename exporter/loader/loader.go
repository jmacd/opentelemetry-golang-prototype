package loader

import (
	"fmt"
	"os"
	"plugin"

	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
)

var (
	pluginName = os.Getenv("OPENTELEMETRY_LIB")
)

func init() {
	if pluginName == "" {
		fmt.Println("Env not set")
		return
	}
	sharedObj, err := plugin.Open(pluginName)
	if err != nil {
		fmt.Println("OPEN failed: ", err)
		return
	}

	obsPlugin, err := sharedObj.Lookup("Observer")
	if err != nil {
		fmt.Println("Symbol not found: ", err)
		return
	}

	obs, ok := obsPlugin.(observer.Observer)
	if !ok {
		fmt.Println("Symbol not an observer: ", err)
		return
	}
	observer.RegisterObserver(obs)
}
