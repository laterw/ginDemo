package commonUtil

import (
	result "flow/entity"
	"flow/redis"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/withlin/canal-go/client"
	pbe "github.com/withlin/canal-go/protocol/entry"
)

//  初始化  阿里 消费端 canal
func InitCanal() {

	// 192.168.199.17 替换成你的canal server的地址
	// example 替换成-e canal.destinations=example 你自己定义的名字
	connector := client.NewSimpleCanalConnector("127.0.0.1", 11111, "", "", "example", 60000, 60*60*1000)
	err := connector.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// https://github.com/alibaba/canal/wiki/AdminGuide
	//mysql 数据解析关注的表，Perl正则表达式.
	//
	//多个正则之间以逗号(,)分隔，转义符需要双斜杠(\\)
	//
	//常见例子：
	//
	//  1.  所有表：.*   or  .*\\..*
	//	2.  canal schema下所有表： canal\\..*
	//	3.  canal下的以canal打头的表：canal\\.canal.*
	//	4.  canal schema下的一张表：canal\\.test1
	//  5.  多个规则组合使用：canal\\..*,mysql.test1,mysql.test2 (逗号分隔)

	err = connector.Subscribe("edt_flowable.*\\..*")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for {

		message, err := connector.Get(100, nil, nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(300 * time.Millisecond)
			//fmt.Println("===没有数据了===")
			continue
		}

		printEntry(message.Entries)

	}
}

func printEntry(entrys []pbe.Entry) {

	for _, entry := range entrys {
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(pbe.RowChange)

		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		checkError(err)

		// 由于没有配置监听库,,,这里 根据 数据库 名称  和  表名称来过滤
		// schema  数据库
		// tableName 表名称
		if rowChange != nil && entry.GetHeader().GetSchemaName() == "schema" && entry.GetHeader().GetTableName() == "tableName" {
			eventType := rowChange.GetEventType()
			header := entry.GetHeader()
			fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))

			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_DELETE {
					printColumn(rowData.GetBeforeColumns(), eventType)
				} else if eventType == pbe.EventType_INSERT {

					printColumn(rowData.GetAfterColumns(), eventType)
				} else {
					fmt.Println("-------> before")
					printColumn(rowData.GetBeforeColumns(), eventType)
					fmt.Println("-------> after")
					printColumn(rowData.GetAfterColumns(), eventType)
				}
			}
		}
	}
}

func printColumn(columns []*pbe.Column, eventType pbe.EventType) {
	for _, col := range columns {
		fmt.Println(fmt.Sprintf("%s : %s  update= %t", col.GetName(), col.GetValue(), col.GetUpdated()))
	}
	result := result.Result{
		ID:        "",
		TaskKey:   "",
		Total:     "",
		Tenant:    "",
		StartUser: "",
		StartTime: "",
		BussKey:   "",
	}
	for i := range columns {
		if columns[i].GetName() == "BUSINESS_KEY_" {
			result.BussKey = columns[i].GetValue()
		}
		if columns[i].GetName() == "START_TIME_" {
			result.StartTime = columns[i].GetValue()
		}
		if columns[i].GetName() == "START_USER_ID_" {
			result.StartUser = redis.Redis.Get(columns[i].GetValue())
		}
	}

	/// 发送到 webSocket
	if result.BussKey != "" && !IsSave(123456) && eventType == pbe.EventType_INSERT {
		SetMessage(123456, result)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
