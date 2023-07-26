package time

import "time"

func NowUnix() int64 {
	return time.Now().Unix()
}

func NowUnixNano() int64 {
	return time.Now().UnixNano()
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func UinxToString(s int64) string {
	return time.Unix(s, 0).Format("2006-01-02 15:04:05")
}

func StringToTime(s string) time.Time {
	t, e := time.Parse("2006-01-02 15:04:05", s)
	if e != nil {
		return time.Time{}
	}
	return t
}
