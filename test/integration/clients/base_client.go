package clients

import (
	"testing"

	core_client "test/integration/clients/core"
)

const Host string = "http://localhost:8383"

func GetClient(t *testing.T) *core_client.Client {
	client, err := core_client.NewClient(Host)
	if err != nil {
		t.Fatal(err)
	}
	return client
}
