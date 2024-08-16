package main

import (
	"go.uber.org/zap/zapcore"
)

type options struct {
	apiServerURL                string
	applicationSecretsNamespace string
	externalAPIPort             int
	kubeConfig                  string
	logLevel                    *zapcore.Level
	proxyCacheTTL               int
	proxyPort                   int
	proxyPortCompass            int
	proxyTimeout                int
	requestTimeout              int
}
