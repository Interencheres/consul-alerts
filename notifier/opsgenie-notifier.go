package notifier

import (
	"fmt"

	alerts "github.com/opsgenie/opsgenie-go-sdk/alerts"
	ogcli "github.com/opsgenie/opsgenie-go-sdk/client"

	log "github.com/Sirupsen/logrus"
)

type OpsGenieNotifier struct {
	ClusterName string
	ApiKey   string
}

func (notifier *OpsGenieNotifier) Notify(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()
	var title string
	text := fmt.Sprintf(header, notifier.ClusterName, overallStatus, fail, warn, pass)

	for _, message := range messages {
		text += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		text += fmt.Sprintf("\n%s", message.Output)
		title += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
	}

	client := new(ogcli.OpsGenieClient)
	client.SetApiKey(notifier.ApiKey)

	alertCli, cliErr := client.Alert()

	if cliErr != nil {
                log.Println("Opsgenie notification trouble with client")
		return false
	}

	// create the alert
	req := alerts.CreateAlertRequest{
		Message:      title,
		Description:  text,
		Source:       "consul",
	}
	response, alertErr := alertCli.Create(req)

	if alertErr != nil {
        log.Println("Opsgenie notification trouble.", response.Status)
		return false
	}

        log.Println("Opsgenie notification send.")
	return true
}
