package snowflake

import (
	"github.com/sony/sonyflake"
	"time"
)

var sf *sonyflake.Sonyflake

func Init() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Date(2021, 7, 28, 0, 0, 0, 0, time.UTC),
	})

	if sf == nil {
		panic("sonyflake not created")
	}
}

func MustID() uint64 {
	id, err := sf.NextID()
	if err == nil {
		return id
	}

	sleep := 1

	for {
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		id, err := sf.NextID()
		if err == nil {
			return id
		}
		sleep *= 2
	}
}
