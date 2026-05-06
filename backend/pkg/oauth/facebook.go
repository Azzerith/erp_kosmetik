package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type FacebookOAuth struct {
	config *oauth2.Config
}

type FacebookUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

func NewFacebookOAuth(clientID, clientSecret, redirectURL string) *FacebookOAuth {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"email",
			"public_profile",
		},
		Endpoint: facebook.Endpoint,
	}
	
	return &FacebookOAuth{config: config}
}

func (f *FacebookOAuth) GetAuthURL(state string) string {
	return f.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (f *FacebookOAuth) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return f.config.Exchange(ctx, code)
}

func (f *FacebookOAuth) GetUserInfo(ctx context.Context, token *oauth2.Token) (*FacebookUserInfo, error) {
	client := f.config.Client(ctx, token)
	url := "https://graph.facebook.com/me?fields=id,name,email,picture"
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var userInfo FacebookUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}
	
	return &userInfo, nil
}

func (f *FacebookOAuth) VerifyAccessToken(accessToken string) (*FacebookUserInfo, error) {
	url := fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=id,name,email,picture", accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var userInfo FacebookUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}
	
	return &userInfo, nil
}