package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"test/integration/clients"
)

func TestMainIndexCheck(t *testing.T) {
	ctx := t.Context()
	cs := setup(ctx)
	defer cleanup(ctx, cs)
	client := clients.GetClient(t)
	response, err := client.GetApiMainIndexHtml(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)

}
