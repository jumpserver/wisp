package recorderstorage

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"

	"github/jumpserver/wisp/pkg/jms-sdk-go/model"
)

type ESCommandStorage struct {
	Hosts   []string
	Index   string
	DocType string

	InsecureSkipVerify bool
}

func (es ESCommandStorage) BulkSave(commands []*model.Command) (err error) {
	var buf bytes.Buffer
	transport := http.DefaultTransport.(*http.Transport).Clone()
	tlsClientConfig := &tls.Config{InsecureSkipVerify: es.InsecureSkipVerify}
	transport.TLSClientConfig = tlsClientConfig
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: es.Hosts,
		Transport: transport,
	})
	if err != nil {
		return err
	}
	for _, item := range commands {
		meta := []byte(fmt.Sprintf(`{ "index" : { } }%s`, "\n"))
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}
		data = append(data, "\n"...)
		buf.Write(meta)
		buf.Write(data)
	}

	response, err := esClient.Bulk(bytes.NewReader(buf.Bytes()),
		esClient.Bulk.WithIndex(es.Index), esClient.Bulk.WithDocumentType(es.DocType))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var (
		blk           *bulkResponse
		numErrors     int64
		errContentBuf bytes.Buffer
	)
	if response.IsError() {
		_, _ = buf.ReadFrom(response.Body)
		return fmt.Errorf("es failure to bulk save: %s", errContentBuf.String())
	}
	if err = json.NewDecoder(response.Body).Decode(&blk); err != nil {
		return err
	}
	for _, d := range blk.Items {
		if d.Index.Status > 201 {
			numErrors++
			errMsg := fmt.Sprintf("%d). [%d]: %s: %s: %s: %s\n",
				numErrors,
				d.Index.Status,
				d.Index.Error.Type,
				d.Index.Error.Reason,
				d.Index.Error.Cause.Type,
				d.Index.Error.Cause.Reason)
			errContentBuf.WriteString(errMsg)
		}
	}
	if numErrors > 0 {
		return fmt.Errorf("es failure to bulk save: %s", errContentBuf.String())
	}
	return nil
}

func (es ESCommandStorage) TypeName() string {
	return "es"
}

// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-bulk.html#bulk-api-response-body
type bulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
}
