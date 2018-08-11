package collector

import (
	"net/http"
	"golang.org/x/oauth2"
	"github.com/sirupsen/logrus"
	"github.com/google/go-github/github"
	"context"
	"github.com/spf13/viper"
)

//NewV3Client returns an authenticated or anonymous GitHub v3 client
func NewV3Client() (client *github.Client, ctx context.Context) {
	ctx = context.Background()
	var tc *http.Client
	// Get authenticadtion if token present
	token := viper.GetString("token")
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
