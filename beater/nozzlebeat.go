package beater

import (
	"fmt"
	"crypto/tls"
	"time"

	"github.com/cloudfoundry-incubator/uaago"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/mikeh-elastic/nozzlebeat/config"
)

type Nozzlebeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Nozzlebeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Nozzlebeat) Run(b *beat.Beat) error {
	logp.Info("nozzlebeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	uaaClient, err := uaago.NewClient(bt.config.FirehoseuaaURL)
	if err != nil {
		logp.Info("Error creating uaa client: %s", err.Error())
	}

	var authToken string
	authToken, err = uaaClient.GetAuthToken(bt.config.Firehoseuser, bt.config.Firehosepassword, true)
	if err != nil {
		logp.Info("Error getting oauth token: %s. Please check your username and password.", err.Error())
	}

	connection := consumer.New(bt.config.FirehosetrafficControllerURL, &tls.Config{InsecureSkipVerify: true}, nil)

	const firehoseSubscriptionId = "firehose-a"
	var (
		msgChan   <-chan *events.Envelope
		errorChan <-chan error
	)
	msgChan, errorChan = connection.Firehose(firehoseSubscriptionId, authToken)

	go func() {
		for err := range errorChan {
			logp.Info("%v", err.Error())
		}
	}()

	for {
	        select {
	        case <-bt.done:
	            return nil
	        case msg := <-msgChan:
			event := beat.Event{
				Timestamp: time.Unix(0, msg.GetTimestamp()),
				Fields: common.MapStr{
					"origin": msg.GetOrigin(),
					"eventType": msg.GetEventType(),
					"deployment": msg.GetDeployment(),
					"job": msg.GetJob(),
					"index": msg.GetIndex(),
					"ip": msg.GetIp(),
					//Not sure if the map[string]string type will go into the []string "tags": msg.GetTags(),
					"httpStartStop": msg.GetHttpStartStop(),
					"logMessage": msg.GetLogMessage(),
					"valueMetric": msg.GetValueMetric(),
					"counterEvent": msg.GetCounterEvent(),
					"error": msg.GetError(),
					"containerMetric": msg.GetContainerMetric(),
				},
			}
			bt.client.Publish(event)
		}
	}
}

func (bt *Nozzlebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
