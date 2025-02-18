package callbacks_test

import (
	"context"
	"github.com/habiliai/habiliai/api/pkg/callbacks"
	habconfig "github.com/habiliai/habiliai/api/pkg/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTweet(t *testing.T) {
	ctx := context.TODO()
	conf := habconfig.ReadHabApiConfig("")

	if conf.Twitter.Validate() != nil {
		t.Skipf("Twitter config is not valid")
	}

	require.NoError(t, callbacks.Tweet(ctx, conf.Twitter, "Hello, world!"))
}
