package models

import (
	"encoding/json"
	"time"
)

// Duration embeds time.Duration and makes it more JSON-friendly.
// Instead of marshaling and unmarshaling as int64 it uses strings, like "5m" or "0.5s".
type Duration time.Duration

func (d Duration) TimeDuration() time.Duration {
	return time.Duration(d)
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	val, err := time.ParseDuration(str)
	*d = Duration(val)
	return err
}

func (d Duration) String() string {
	return time.Duration(d).String()
}
