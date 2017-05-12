package app

import (
	"cn/com/hengwei/commons"
	"cn/com/hengwei/commons/httputils"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/revel/revel"
	sso "github.com/three-plus-three/sso/client"
	"github.com/three-plus-three/sso/client/revel_sso"
)

func initSSO(env *commons.Environment) revel_sso.CheckFunc {
	ssoURL := env.GetServiceConfig(commons.ENV_WSERVER_PROXY_ID).UrlFor(env.DaemonUrlPath, "/sso")
	ssoClient, err := sso.NewClient(ssoURL)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
		return nil
	}

	return revel_sso.SSO(ssoClient, 30*time.Minute, func(req *http.Request) url.URL {
		copyURL := *req.URL
		copyURL.Scheme = ""
		copyURL.Host = ""
		copyURL.Path = httputils.JoinURLPath(revel.AppRoot, copyURL.Path)
		return copyURL
	})
}
