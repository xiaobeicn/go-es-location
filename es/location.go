package es

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/xiaobeicn/go-es-location/global"
)

type Location struct {
	IP           string            `json:"ip,omitempty"`
	City         string            `json:"city,omitempty"`
	CountryLong  string            `json:"country_long,omitempty"`
	CountryShort string            `json:"country_short,omitempty"`
	Region       string            `json:"region,omitempty"`
	Timezone     string            `json:"timezone,omitempty"`
	Zipcode      string            `json:"zipcode,omitempty"`
	Latitude     float64           `json:"latitude,omitempty"`
	Longitude    float64           `json:"longitude,omitempty"`
	Location     *elastic.GeoPoint `json:"location,omitempty"`
}

func IndexName() string {
	return "location_v1"
}

func setting() string {
	return `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"city":{
					"type":"keyword"
				},
				"country_long":{
					"type":"keyword"
				},
				"country_short":{
					"type":"keyword"
				},
				"region":{
					"type":"keyword"
				},
				"timezone":{
					"type":"keyword"
				},
				"zipcode":{
					"type":"keyword"
				},
				"ip":{
					"type": "ip"
				},
				"longitude":{
					"type": "float"
				},
				"latitude":{
					"type": "float"
				},
				"location":{
					"type":"geo_point"
				}
			}
		}
	}`
}

func CreateIndex(ctx *gin.Context) (string, error) {
	name := IndexName()
	exists, err := global.GElastic.IndexExists(name).Do(ctx)
	if err != nil {
		return "", err
	}
	if !exists {
		createIndex, err := global.GElastic.CreateIndex(name).BodyString(setting()).Do(ctx)
		if err != nil {
			return "", err
		}
		if !createIndex.Acknowledged {
			return "", errors.New("Un acknowledged")
		}
	}
	return name, nil
}

func Index(ctx *gin.Context, body *Location) (*elastic.IndexResponse, error) {
	return global.GElastic.Index().Index(IndexName()).BodyJson(body).Do(ctx)
}

func Search(ctx *gin.Context, lat, lon float64) (total int64, list []Location, err error) {
	q := elastic.NewGeoDistanceQuery("location")
	q = q.GeoPoint(elastic.GeoPointFromLatLon(lat, lon))
	q = q.Distance("3km")
	searchResult, err := global.GElastic.Search(IndexName()).Query(q).Do(ctx)
	if err != nil {
		return
	}
	total = searchResult.TotalHits()
	if total > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var t Location
			err := json.Unmarshal(hit.Source, &t)
			if err == nil {
				list = append(list, t)
			}
		}
	}
	return
}
