package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
	"github.com/xiaobeicn/go-es-location/es"
	"github.com/xiaobeicn/go-es-location/util/location"
	"github.com/xiaobeicn/go-es-location/util/net"
)

func HomeHandler(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		ip = net.RemoteIp(c.Request)
	}
	results, err := location.SearchAll(ip)
	if err != nil {
		c.JSON(500, gin.H{
			"ip":  ip,
			"err": err.Error(),
		})
	}

	lat := cast.ToFloat64(results.Latitude)
	lon := cast.ToFloat64(results.Longitude)
	rep, err := es.Index(c, &es.Location{
		IP:           ip,
		City:         results.City,
		CountryLong:  results.Country_long,
		CountryShort: results.Country_short,
		Region:       results.Region,
		Timezone:     results.Timezone,
		Zipcode:      results.Zipcode,
		Latitude:     lat,
		Longitude:    lon,
		Location:     elastic.GeoPointFromLatLon(lat, lon),
	})
	if err != nil {
		c.JSON(500, gin.H{
			"ip":  ip,
			"err": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"ip":     ip,
			"status": rep.Status,
			"id":     rep.Id,
			"result": rep.Result,
		})
	}
}

func NearHandler(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		ip = net.RemoteIp(c.Request)
	}
	locationResults, err := location.SearchAll(ip)
	if err != nil {
		c.JSON(500, gin.H{
			"ip":  ip,
			"err": err.Error(),
		})
	}
	total, list, err := es.Search(c, cast.ToFloat64(locationResults.Latitude), cast.ToFloat64(locationResults.Longitude))
	if err != nil {
		c.JSON(500, gin.H{
			"ip":  ip,
			"err": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"ip":    ip,
		"total": total,
		"list":  list,
	})
}

func IpHandler(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		ip = net.RemoteIp(c.Request)
	}
	results, err := location.SearchAll(ip)
	if err != nil {
		c.JSON(500, gin.H{
			"ip":  ip,
			"err": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"ip":            ip,
		"country_short": results.Country_short,
		"country_long":  results.Country_long,
		"region":        results.Region,
		"city":          results.City,
		"latitude":      results.Latitude,
		"longitude":     results.Longitude,
		"zipcode":       results.Zipcode,
		"timezone":      results.Timezone,
	})
}
