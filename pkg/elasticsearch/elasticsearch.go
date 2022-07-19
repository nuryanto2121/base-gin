package elasticsearch

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
	"gitlab.com/369-engineer/369backend/account/pkg/setting"
)

var Client *elastic.Client

func Setup() {
	// Instantiate a client instance of the elastic library
	var err error
	options := make([]elastic.ClientOptionFunc, 0)

	options = append(options, elastic.SetURL(setting.ElasticSearchSetting.ElasticHost))
	options = append(options, elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)))
	options = append(options, elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	options = append(options, elastic.SetTraceLog(log.New(os.Stderr, "", 0)))

	isAuthEnable := setting.ElasticSearchSetting.ElasticAuth

	if isAuthEnable {
		options = append(options, elastic.SetSniff(false))
		options = append(options, elastic.SetHealthcheck(false))
		options = append(options, elastic.SetBasicAuth(setting.ElasticSearchSetting.ElasticAuthUsername, setting.ElasticSearchSetting.ElasticAuthPassword))
	} else {
		options = append(options, elastic.SetSniff(true))
		options = append(options, elastic.SetHealthcheckInterval(5*time.Second))
	}
	Client, err = elastic.NewClient(options...)

	if err != nil {
		// (Bad Request): Failed to parse content to map if mapping bad
		log.Panic("elastic.NewClient() ERROR: ", err)
	}
	go CheckOrCreateLoggingIndices()

	log.Println("elastic client:", Client)

}

func CheckOrCreateLoggingIndices() (err error) {
	mappings := `
	{
		"settings": {
			"number_of_shards": 2,
			"number_of_replicas": 1
		},
		"mappings": {
			"properties": {
				"audit_log_id int": {
					"type": "integer"
				},
				"created_at date": {
					"type": "date"
				},
				"updated_at date": {
					"type": "date"
				},
				"level str": {
					"type": "text"
				},
				"uuid str": {
					"type": "text"
				},
				"func_name str": {
					"type": "text"
				},
				"file_name str": {
					"type": "text"
				},
				"line int": {
					"type": "integer"
				},
				"time str": {
					"type": "text"
				},
				"message str": {
					"type": "text"
				}
			}
		}
	}
	
	`
	ctx := context.Background()
	indexExist, err := Client.IndexExists("audit_log_chat").Do(ctx)
	if err != nil {
		return err
	}

	if !indexExist {
		_, err = Client.CreateIndex("audit_log_chat").
			Body(mappings).
			Do(ctx)
	}

	return err
}
