package gron

import (
	"fmt"
	"testing"
	"time"

	"github.com/roylee0704/gron/xtime"
)

func TestPeriodicNext(t *testing.T) {

	tests := []struct {
		time   string
		period time.Duration
		want   string
	}{

		//Wraps around days
		{
			time:   "Mon Jun 6 23:46 2016",
			period: 15 * time.Minute,
			want:   "Mon Jun 7 00:01 2016",
		},
		{
			time:   "Mon Jun 6 23:59:40 2016",
			period: 21 * time.Second,
			want:   "Mon Jun 7 00:00:01 2016",
		},
		{
			time:   "Mon Jun 6 23:30:19 2016",
			period: 29*time.Minute + 41*time.Second,
			want:   "Mon Jun 7 00:00:00 2016",
		},
		{
			time:   "Mon Jun 6 23:46:20 2016",
			period: 15*time.Minute + 40*time.Second,
			want:   "Mon Jun 7 00:02:00 2016",
		},

		// Wrap around months
		{
			time:   "Mon Jun 6 16:49 2016",
			period: 30 * xtime.Day, // adds 30 days, equates to #days in June.
			want:   "Mon Jul 6 16:49 2016",
		},

		// Wrap around minute, hour, day, month, and year
		{

			time:   "Sat Dec 31 23:59:59 2016",
			period: 1 * time.Second,
			want:   "Sun Jan 1 00:00:00 2017",
		},

		// Round to nearest second on the period
		{
			time:   "Mon Jun 6 12:45 2016",
			period: 15*time.Minute + 100*time.Nanosecond,
			want:   "Mon Jun 6 13:00 2016",
		},

		// Round to 1 second if the duration is less
		{
			time:   "Mon Jun 6 17:38:01 2016",
			period: 59 * time.Millisecond,
			want:   "Mon Jun 6 17:38:02 2016",
		},

		// Round to nearest second when calculating the next time.
		{"Mon Jun 6 17:38:01.009 2016", 15 * time.Minute, "Mon Jun 6 17:53:01 2016"},

		// Round to nearest second for both
		{"Mon Jun 6 17:38:01.009 2016", 15*time.Minute + 50*time.Nanosecond, "Mon Jun 6 17:53:01 2016"},
	}

	for i, test := range tests {

		got := Every(test.period).Next(getTime(test.time))
		want := getTime(test.want)

		fmt.Println(got)
		fmt.Println(want)

		if got != want {
			t.Errorf("case[%d], %s, \"%s\": (want) %v != %v (got)", i, test.time, test.period, want, got)
		}
	}

}

func getTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	t, err := time.Parse("Mon Jan 2 15:04 2006", value)
	if err != nil {
		t, err = time.Parse("Mon Jan 2 15:04:05 2006", value)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05-0700", value)
			if err != nil {
				panic(err)
			}
			// Daylight savings time tests require location
			if ny, err := time.LoadLocation("UTC"); err == nil {
				t = t.In(ny)
			}
		}
	}

	return t
}