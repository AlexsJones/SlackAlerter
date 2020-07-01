package subscriptions

import (
	"encoding/json"
	"fmt"
	"github.com/AlexsJones/kubeops/lib/subscription"
	"github.com/slack-go/slack"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"
	"os"
	"strconv"
	"time"
)

type SlackAlertPod struct{}

func (SlackAlertPod) WithElectedResource() interface{} {

	return &v1.Pod{}
}

func (SlackAlertPod) WithEventType() []watch.EventType {

	return []watch.EventType {watch.Deleted}
}

func (SlackAlertPod) OnEvent(msg subscription.Message) {

	pod := msg.Event.Object.(*v1.Pod)
	klog.Infof("Pod deleted %s",pod.Name)
	tok := os.Getenv("SLACK_WEBHOOK")
	if tok != ""  {
		attachment := slack.Attachment{
			Color:         "good",
			Fallback:      "You successfully posted by Incoming Webhook URL!",
			AuthorName:    "slackalerter",
			AuthorSubname: "github.com",
			AuthorLink:    "https://github.com/slack-go/slack",
			AuthorIcon:    "https://avatars2.githubusercontent.com/u/652790",
			Text:          fmt.Sprintf("<!channel> Pod has been deleted %s",pod.Name),
			Footer:        "slackalerter",
			FooterIcon:    "https://platform.slack-edge.com/img/default_application_icon.png",
			Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
		}
		msg := slack.WebhookMessage{
			Attachments: []slack.Attachment{attachment},
		}

		err := slack.PostWebhook(tok, &msg)
		if err != nil {
			fmt.Println(err)
		}
	}else {
		klog.Info("no slack token set, ignoring...")
	}
}


