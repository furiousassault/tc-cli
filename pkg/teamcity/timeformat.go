package teamcity

import "time"

const timeFormat = "20060102T150405+0000"

func ParseTCTimeFormat(tcTime string) (timestamp time.Time, err error) {
	return time.Parse(timeFormat, tcTime)
}
