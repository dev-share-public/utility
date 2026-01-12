package timezo

import (
	"fmt"
	"time"
)

func InitTime(timezone string) error {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	time.Local = loc
	fmt.Println("✅ Time Zone Set To : ", time.Local.String())
	return nil
}
