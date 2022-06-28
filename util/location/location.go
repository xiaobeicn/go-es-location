package location

import (
	"github.com/ip2location/ip2location-go"
)

func SearchAll(ip string) (res ip2location.IP2Locationrecord, err error) {
	db, err := ip2location.OpenDB("./data/IP2LOCATION-LITE-DB11.BIN")
	if err != nil {
		return
	}
	res, err = db.Get_all(ip)
	db.Close()
	return
}
