package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/require"
)

func TestFormatUserLogsRemovesAdminOnlyFields(t *testing.T) {
	logs := []*Log{
		{
			Id:          99,
			ChannelName: "private-channel",
			Other: common.MapToJsonStr(map[string]interface{}{
				"admin_info":    "debug",
				"reject_reason": "internal policy detail",
				"stream_status": "internal stream detail",
				"public":        "visible",
			}),
		},
	}

	formatUserLogs(logs, 4)

	require.Equal(t, 5, logs[0].Id)
	require.Empty(t, logs[0].ChannelName)

	otherMap, err := common.StrToMap(logs[0].Other)
	require.NoError(t, err)
	require.Equal(t, "visible", otherMap["public"])
	require.NotContains(t, otherMap, "admin_info")
	require.NotContains(t, otherMap, "reject_reason")
	require.NotContains(t, otherMap, "stream_status")
}
