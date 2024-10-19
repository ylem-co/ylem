package services

import (
	"io"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"net/http"
	"net/url"
	"ylem_integrations/config"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type SlackAuthorizedGrant struct {
	Ok          bool   `json:"ok"`
	Error       string `json:"error,omitempty"`
	Scopes      string `json:"scope,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	BotUserId   string `json:"bot_user_id,omitempty"`
}

const SlackOAuth2AuthorizeUrl = "https://slack.com/api/oauth.v2.access"

func SlackGrantAuthorization(Token string) (*SlackAuthorizedGrant, error) {
	var Config config.Config
	err := envconfig.Process("", &Config)
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	Data := url.Values{}
	Data.Set("client_id", Config.Slack.ClientId)
	Data.Set("client_secret", Config.Slack.ClientSecret)
	Data.Set("code", Token)

	Request, err := http.NewRequest("POST", SlackOAuth2AuthorizeUrl, strings.NewReader(Data.Encode()))

	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	Request.Header.Set("Content-Length", strconv.Itoa(len(Data.Encode())))

	client := &http.Client{} // that looks fishy
	resp, err := client.Do(Request)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	var Grant SlackAuthorizedGrant
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &Grant)
	if err != nil {
		log.Error(err.Error())

		return nil, err
	}

	if !Grant.Ok {
		message := "Slack access token error: " + Grant.Error
		log.Error(message)

		return nil, errors.New(message)
	}

	return &Grant, nil
}

func CreateSlackGrantLink(State string) string {
	return fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&scope=channels:join,chat:write,channels:read,groups:write,groups:read&state=%s",
		config.Cfg().Slack.ClientId,
		State,
	)
}

func GetSlackChannelIdFromName(SlackClient *slack.Client, ChannelName string) (*string, error) {
	configuration := slack.GetConversationsParameters{
		ExcludeArchived: true,
		Limit:           150,
		Types:           []string{"private_channel", "public_channel"},
	}

	for ok := true; ok; ok = configuration.Cursor != "" {
		channels, cursor, err := SlackClient.GetConversations(&configuration)

		if err != nil {
			log.Error(err.Error())

			return nil, err
		}

		for _, v := range channels {
			if v.Name == ChannelName {
				return &v.ID, nil
			}
		}

		configuration.Cursor = cursor
	}

	return nil, nil
}

func JoinSlackChannel(SlackClient *slack.Client, ChannelId string) error {
	_, _, _, err := SlackClient.JoinConversation(ChannelId)

	if err != nil {
		log.Error(err.Error())

		return nil
	}

	return nil
}
