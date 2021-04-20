package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/sunker/redshift-datasource/pkg/redshift"
)



func main() {
	// Start listening to requests send from Grafana. This call is blocking so
	// it wont finish until Grafana shutsdown the process or the plugin choose
	// to exit close down by itself
	// err := datasource.Serve(newDatasource())
	ds := redshift.NewDatasource()
	err := experimental.DoGRPC("redshift", ds)


	// Log any error if we could start the plugin.
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
