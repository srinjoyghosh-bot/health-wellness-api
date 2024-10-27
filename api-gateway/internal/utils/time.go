package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ProtoTimestampToTime(ts *timestamppb.Timestamp) time.Time {
	return ts.AsTime()
}

func TimeToProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
