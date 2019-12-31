package todo

import (
	"fmt"
	"testing"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var InfluxClient client.Client

func InitInfluxDB(influxUrl, userName, pwd, databaseName string) {
	InfluxClient, _ = initInflux(influxUrl, userName, pwd, databaseName)
}

func initInflux(host, userName, pwd, databaseName string) (client.Client, error) {
	InfluxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://%s", host),
		Username: userName,
		Password: pwd,
		Timeout:  time.Second * 5,
	})
	if err != nil {
		return InfluxClient, err
	}
	// 建库
	createDbSQL := client.NewQuery(fmt.Sprintf("CREATE DATABASE %s", databaseName), "", "")

	createCQ2m := client.NewQuery(fmt.Sprintf("CREATE CONTINUOUS QUERY cq_exchanger_2m ON %s "+
		"BEGIN SELECT mean(value) AS value INTO raw_data_2m FROM raw_data GROUP BY unique_id, time(2m) END", databaseName), databaseName, "")
	createCQ1h := client.NewQuery(fmt.Sprintf("CREATE CONTINUOUS QUERY cq_exchanger_1h ON %s "+
		"BEGIN SELECT mean(value) AS value INTO raw_data_1h FROM raw_data GROUP BY unique_id, time(1h) END", databaseName), databaseName, "")
	// 过期策略
	createRPSQL := client.NewQuery(fmt.Sprintf("CREATE RETENTION POLICY default ON %s DURATION 3d REPLICATION 1 DEFAULT", databaseName), databaseName, "")

	if _, err := InfluxClient.Query(createDbSQL); err != nil {
		return InfluxClient, err
	}
	if _, err := InfluxClient.Query(createCQ2m); err != nil {
		return InfluxClient, err
	}
	if _, err := InfluxClient.Query(createCQ1h); err != nil {
		return InfluxClient, err
	}
	if _, err := InfluxClient.Query(createRPSQL); err != nil {
		return InfluxClient, err
	}

	return InfluxClient, nil
}

func TestInfluxDB(t *testing.T) {
	databaseName := "appXXX"
	tableName := "tableName"
	InitInfluxDB("127.0.0.1:8086", "root", "root", databaseName)
	now := time.Now()

	//批量插入
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "ms",
		Database:  databaseName,
	})

	point, _ := client.NewPoint(
		tableName,
		map[string]string{"tag1": "tag1V"}, //	tag 唯一索引
		map[string]interface{}{"fieldName": "fieldNameValue"}, now)
	bp.AddPoint(point)

	fmt.Println(InfluxClient)
	err := InfluxClient.Write(bp)
	fmt.Println(err)

	//	查询
	query := fmt.Sprintf(`select * from "%s"`, tableName)
	resp, err := InfluxClient.Query(client.Query{Command: query, Database: databaseName})
	if err != nil || len(resp.Err) != 0 {
		fmt.Println(err, "->", resp.Err)
	}
	fmt.Println(resp.Results)
}
