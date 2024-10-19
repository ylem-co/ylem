package slack

import (
	"fmt"
	
	messaging "github.com/ylem-co/shared-messaging"
	"github.com/slack-go/slack"
)

func SendSlackMessage(ChannelId string, Title string, Text string, Severity string, AccessToken string) error {
	api := slack.New(AccessToken)

	severityToColor := map[string]string {
		messaging.TaskSeverityCritical: "#8b0000",
		messaging.TaskSeverityHigh: "#8b0000",
		messaging.TaskSeverityMedium: "#DEC20B",
		messaging.TaskSeverityLowest: "#006400",
		messaging.TaskSeverityLow: "#006400",
	}

	color, ok := severityToColor[Severity]
	if !ok {
		return fmt.Errorf(`unknown task severity "%s"`, Severity)
	}

	headerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s*", Title), false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	attachment := slack.Attachment{
		Text: Text,
		Color: color,
	}

	_, _, err := api.PostMessage(
		ChannelId,
		slack.MsgOptionBlocks(headerSection),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)

	return err
}
