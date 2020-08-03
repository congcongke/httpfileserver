package server

import (
	"fmt"
	"net/http"

	"github.com/congcongke/httpfileserver/pkg/config"
)

func Run(conf *config.Config) {
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), LoadFromConfig(conf))
	if err != nil {
		panic(err)
	}
}
