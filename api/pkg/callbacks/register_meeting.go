package callbacks

import "context"

func RegisterMeeting(s *service, ctx context.Context, args []byte, metadata Metadata) (any, error) {
	return nil, nil
}

func init() {
	dispatchFunctions["register_meeting"] = RegisterMeeting
}
