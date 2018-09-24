package collector

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thales-e-security/contribstats/pkg/config"
	"golang.org/x/oauth2"
	"net/http"
)

//NewV3Client returns an authenticated or anonymous GitHub v3 client
func NewV3Client(constants config.Config) (client *github.Client, ctx context.Context) {
	ctx = context.Background()
	var tc *http.Client
	// Get authenticadtion if token present
	token := constants.Token
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctx, ts)
	} else {
		logrus.Warn("No token provided, you are not likely to get much details as most organizations default to private membership")
		logrus.Warnf("Try adding token to your config at: %v", viper.ConfigFileUsed())
	}
	client = github.NewClient(tc)
	return
}
