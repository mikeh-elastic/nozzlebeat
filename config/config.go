// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Firehoseuser string `config:"firehoseuser"`
	Firehosepassword string `config:"firehosepassword"`
	FirehosetrafficControllerURL string `config:"firehosetrafficControllerURL"`
	FirehoseuaaURL string `config:"firehoseuaaURL"`
}

var DefaultConfig = Config{
	Firehoseuser: "example-nozzle",
	Firehosepassword: "example-nozzle",
	FirehosetrafficControllerURL: "wss://doppler.bosh-lite.com:443",
	FirehoseuaaURL: "https://uaa.bosh-lite.com",
}
