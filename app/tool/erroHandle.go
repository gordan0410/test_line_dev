package tool

import "github.com/rs/zerolog/log"

func ErroHandle(packageName, serviceName, functionName string, err error) {
	log.Error().Caller().Err(err).Msgf("packageName: %s, svcName: %s,funcName: %s", packageName, serviceName, functionName)
}
