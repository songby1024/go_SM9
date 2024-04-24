package createData

import (
	"math/rand"
	"strconv"
	"time"
)

func CreateRandem() string {
	rand.Seed(time.Now().UnixNano())
	res := strconv.FormatInt(rand.Int63n(9999999999999999), 10)
	if len(res) < 16 {
		for {
			res += "0"
			if len(res) == 16 {
				break
			}
		}
	}
	return res
}
