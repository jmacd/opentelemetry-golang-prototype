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
		return
	}
	sharedObj, err := plugin.Open(pluginName)
	if err != nil {
		fmt.Println("Open failed", pluginName, err)
		return
	}

	obsPlugin, err := sharedObj.Lookup("Observer")
	if err != nil {
		fmt.Println("Observer not found", pluginName, err)
		return
	}

	obs, ok := obsPlugin.(*observer.Observer)
	if !ok {
		fmt.Printf("Observer not valid\n")
		return
	}
	observer.RegisterObserver(*obs)
}
