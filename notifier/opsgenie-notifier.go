package notifier

import (
	"fmt"
	"net/url"

    alerts "github.com/opsgenie/opsgenie-go-sdk/alerts"
    ogcli "github.com/opsgenie/opsgenie-go-sdk/client"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type OpsGenieNotifier struct {
	ClusterName string
	RoomId      string
	AuthToken   string
	BaseURL     string
}

func (notifier *OpsGenieNotifier) Notify(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	text := fmt.Sprintf(header, notifier.ClusterName, overallStatus, fail, warn, pass)

	for _, message := range messages {
		text += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		text += fmt.Sprintf("\n%s", message.Output)
	}

	level := "green"
	if fail > 0 {
		level = "red"
	} else if warn > 0 {
		level = "yellow"
	}

    client := new (ogcli.OpsGenieClient)
    client.SetApiKey(notifier.ApiKey)

    alertCli, cliErr := client.Alert()

    if cliErr != nil {
		log.Printf("Error instanciating OpsGenie's client: %s\n", cliErr)
		return false
    }

    // create the alert
    req := alerts.CreateAlertRequest{
        Message: "appserver1 down",
        Description: "cpu usage is over 60%",
        Source: "consul",
        Entity: notifier.ClusterName,
        Actions: []string{"ping", "restart"},
        Tags: []string{"network", "operations"},
        // XXX needed ?
        Recipients: []string{"john.smith@acme.com", "admin@acme.com"}
    }
    response, alertErr := alertCli.Create(req)

    if alertErr != nil {
		log.Printf("Error sending notification to OpsGenie: %s\n", alertErr)
		log.Printf("Server returns %+v\n", response)
		return false
	}

	return true
}
