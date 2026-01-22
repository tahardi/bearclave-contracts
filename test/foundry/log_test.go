package foundry_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tahardi/bearchain/test/foundry"
)

//go:embed testdata/log.json
var logJSON []byte

func TestLog_JSON(t *testing.T) {
	t.Run("happy path - round trip", func(t *testing.T) {
		// given
		want := logJSON

		// when
		log := &foundry.Log{}
		err := log.UnmarshalJSON(want)
		require.NoError(t, err)

		got, err := log.MarshalJSON()

		// then
		require.NoError(t, err)
		require.JSONEq(t, string(want), string(got))
	})
}
