package prometheus

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(serverAddr string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		if err := http.ListenAndServe(serverAddr, mux); err != nil {
			fmt.Printf("Error starting prometheus server: %s\n", err.Error())
		}
	}()
}
