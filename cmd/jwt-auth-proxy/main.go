package main

import (
	"os"

	"github.com/janivihervas/authproxy/internal/server"

	"github.com/janivihervas/authproxy/upstream"
)

type Config struct {
	AzureClientID     string `envconfig:"AZURE_AD_CLIENT_ID"`
	AzureClientSecret string `envconfig:"AZURE_AD_CLIENT_SECRET"`
	AzureTenant       string `envconfig:"AZURE_AD_TENANT"`
}

func main() {
	//var config Config
	//
	//_ = gotenv.OverLoad(".env")
	//err := envconfig.Process("", &config)
	//if err != nil {
	//	panic(err)
	//}
	//
	//m := jwt_auth_proxy.NewMiddleware(&oauth2.Config{
	//	ClientID:     config.AzureClientID,
	//	ClientSecret: config.AzureClientSecret,
	//	Endpoint:     azure.Endpoint(config.AzureTenant),
	//	RedirectURL:  "http://localhost:3000/oauth2/callback",
	//	Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
	//},
	//	memory.New(),
	//	upstream.Echo{},
	//)
	//
	//err = server.RunHTTP(os.Getenv("PORT"), m)
	//if err != nil {
	//	panic(err)
	//}

	err := server.RunHTTP(os.Getenv("PORT"), upstream.Echo{})
	if err != nil {
		panic(err)
	}
}
