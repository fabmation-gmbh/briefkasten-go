package oauth

import (
	"log"
	"net/url"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/battlenet"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/eveonline"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/kakao"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/line"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/mastodon"
	"github.com/markbates/goth/providers/nextcloud"
	"github.com/markbates/goth/providers/okta"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/seatalk"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/yammer"
	"go.uber.org/zap"

	"github.com/fabmation-gmbh/briefkasten-go/internal/config"
)

// Init initializes all required objects.
// If initialization fails, the function panics.
func Init() {
	goth.UseProviders(
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"),
			genCallbackURL("facebook")),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"),
			genCallbackURL("google")),
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), genCallbackURL("gplus")),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), genCallbackURL("github")),
		spotify.New(os.Getenv("SPOTIFY_KEY"), os.Getenv("SPOTIFY_SECRET"), genCallbackURL("spotify")),
		linkedin.New(os.Getenv("LINKEDIN_KEY"), os.Getenv("LINKEDIN_SECRET"), genCallbackURL("linkedin")),
		line.New(os.Getenv("LINE_KEY"), os.Getenv("LINE_SECRET"), genCallbackURL("line"), "profile", "openid", "email"),
		lastfm.New(os.Getenv("LASTFM_KEY"), os.Getenv("LASTFM_SECRET"), genCallbackURL("lastfm")),
		twitch.New(os.Getenv("TWITCH_KEY"), os.Getenv("TWITCH_SECRET"), genCallbackURL("twitch")),
		dropbox.New(os.Getenv("DROPBOX_KEY"), os.Getenv("DROPBOX_SECRET"), genCallbackURL("dropbox")),
		digitalocean.New(os.Getenv("DIGITALOCEAN_KEY"), os.Getenv("DIGITALOCEAN_SECRET"), genCallbackURL("digitalocean"), "read"),
		bitbucket.New(os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"), genCallbackURL("bitbucket")),
		instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), genCallbackURL("instagram")),
		intercom.New(os.Getenv("INTERCOM_KEY"), os.Getenv("INTERCOM_SECRET"), genCallbackURL("intercom")),
		box.New(os.Getenv("BOX_KEY"), os.Getenv("BOX_SECRET"), genCallbackURL("box")),
		salesforce.New(os.Getenv("SALESFORCE_KEY"), os.Getenv("SALESFORCE_SECRET"), genCallbackURL("salesforce")),
		seatalk.New(os.Getenv("SEATALK_KEY"), os.Getenv("SEATALK_SECRET"), genCallbackURL("seatalk")),
		amazon.New(os.Getenv("AMAZON_KEY"), os.Getenv("AMAZON_SECRET"), genCallbackURL("amazon")),
		yammer.New(os.Getenv("YAMMER_KEY"), os.Getenv("YAMMER_SECRET"), genCallbackURL("yammer")),
		onedrive.New(os.Getenv("ONEDRIVE_KEY"), os.Getenv("ONEDRIVE_SECRET"), genCallbackURL("onedrive")),
		azuread.New(os.Getenv("AZUREAD_KEY"), os.Getenv("AZUREAD_SECRET"), genCallbackURL("azuread"), nil),
		battlenet.New(os.Getenv("BATTLENET_KEY"), os.Getenv("BATTLENET_SECRET"), genCallbackURL("battlenet")),
		eveonline.New(os.Getenv("EVEONLINE_KEY"), os.Getenv("EVEONLINE_SECRET"), genCallbackURL("eveonline")),
		kakao.New(os.Getenv("KAKAO_KEY"), os.Getenv("KAKAO_SECRET"), genCallbackURL("kakao")),

		// Auth0 allocates domain per customer, a domain must be provided for auth0 to work
		auth0.New(os.Getenv("AUTH0_KEY"), os.Getenv("AUTH0_SECRET"), genCallbackURL("auth0"), os.Getenv("AUTH0_DOMAIN")),
		nextcloud.NewCustomisedDNS(os.Getenv("NEXTCLOUD_KEY"), os.Getenv("NEXTCLOUD_SECRET"), genCallbackURL("nextcloud"), os.Getenv("NEXTCLOUD_URL")),
		gitea.New(os.Getenv("GITEA_KEY"), os.Getenv("GITEA_SECRET"), genCallbackURL("gitea")),
		okta.New(os.Getenv("OKTA_ID"), os.Getenv("OKTA_SECRET"), os.Getenv("OKTA_ORG_URL"), genCallbackURL("okta"), "openid", "profile", "email"),
		mastodon.New(os.Getenv("MASTODON_KEY"), os.Getenv("MASTODON_SECRET"), genCallbackURL("mastodon"), "read:accounts"),
	)

	// OpenID Connect is based on OpenID Connect Auto Discovery URL (https://openid.net/specs/openid-connect-discovery-1_0-17.html)
	// because the OpenID Connect provider initialize itself in the New(), it can return an error which should be handled or ignored
	// ignore the error for now
	openidConnect, _ := openidConnect.New(os.Getenv("OPENID_CONNECT_KEY"), os.Getenv("OPENID_CONNECT_SECRET"),
		genCallbackURL("openid-connect"), os.Getenv("OPENID_CONNECT_DISCOVERY_URL"))
	if openidConnect != nil {
		goth.UseProviders(openidConnect)
	}
}

// genCallbackURL returns a valid callback url for the provider.
func genCallbackURL(provider string) string {
	u, err := url.JoinPath(config.C.OAuth2.Endpoint, "/api/v1/auth/", provider, "/callback")
	if err != nil {
		log.Fatal("Unable to generate provider URL", zap.Error(err))
	}

	return u
}
