package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func RecurringTimer(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func RandomFromArray(array []string) string {
	return array[rand.Intn(len(array))]
}

func StripPluginCommand(str string, prefix string, plugin string) string {
	return strings.Replace(str, prefix+plugin+" ", "", -1)
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
