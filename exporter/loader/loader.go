package loader

import (
	"fmt"
	"os"
	"plugin"
	"time"

	"github.com/lightstep/opentelemetry-golang-prototype/exporter/observer"
)

func init() {
	pluginName := os.Getenv("OPENTELEMETRY_LIB")
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

func Flush() {
	// TODO add Flush support in loader
	// TODO implement for exporter/{stdout,stderr,buffer}
	time.Sleep(1 * time.Second)
}
