package callbacks

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type MeetingRecords struct {
	Name     string `json:"meeting_name"`
	DateTime string `json:"date_time"`
	Location string `json:"location"`
}

func RegisterMeeting(ctx *Context, args []byte) (any, error) {
	var records MeetingRecords
	if err := json.Unmarshal(args, &records); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal args")
	}

	if err := ctx.UpdateMemory(map[string]any{
		"meeting": records,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func init() {
	dispatchFunctions["register_meeting"] = RegisterMeeting
}
