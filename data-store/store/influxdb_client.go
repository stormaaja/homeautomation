package store

import (
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDBClient struct {
	Client       influxdb2.Client
	Data         []map[string]interface{}
	Bucket       string
	Organization string
}

func NewInfluxDBClient() *InfluxDBClient {
	url := os.Getenv("INFLUXDB_URL")
	bucket := os.Getenv("INFLUXDB_DATABASE")
	token := os.Getenv("INFLUXDB_TOKEN")
	if url == "" || bucket == "" || token == "" {
		log.Println("InfluxDB environment variables not set. Not creating InfluxDB client.")
		return nil
	}
	return &InfluxDBClient{
		Client: influxdb2.NewClientWithOptions(
			url,
			token,
			influxdb2.DefaultOptions().SetBatchSize(20),
		),
		Data:         []map[string]interface{}{},
		Bucket:       bucket,
		Organization: os.Getenv("INFLUXDB_URL_ORGANIZATION"),
	}
}

func (c *InfluxDBClient) AppendItem(
	measurement string,
	location string,
	field string,
	value float64,
) {
	data := map[string]interface{}{
		"measurement": measurement,
		"location":    location,
		"field":       field,
		"value":       value,
		"time":        time.Now(),
	}

	c.Data = append(c.Data, data)
}

func (i *InfluxDBClient) Flush() {
	writeAPI := i.Client.WriteAPI(i.Organization, i.Bucket)

	for _, item := range i.Data {
		p := influxdb2.NewPointWithMeasurement(item["measurement"].(string)).
			SetTime(time.Now()).
			AddTag("location", item["location"].(string)).
			AddField(item["field"].(string), item["value"].(float64))
		writeAPI.WritePoint(p)
	}

	writeAPI.Flush()
	i.Client.Close()
	i.Data = []map[string]interface{}{}
}
