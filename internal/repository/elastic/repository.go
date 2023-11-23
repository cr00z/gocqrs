package elastic

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cr00z/gocqrs/internal/domain"
	"github.com/olivere/elastic"
)

type ElasticRepository struct {
	client *elastic.Client
}

func NewElasticRepository(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{
		client: client,
	}, nil
}

func (r *ElasticRepository) Close() {
	// not implemented
}

func (r *ElasticRepository) Insert(ctx context.Context, message domain.Message) error {
	_, err := r.client.Index().
		Index("messages").
		Type("message").
		Id(message.ID).
		BodyJson(message).
		Refresh("wait_for").Do(ctx)
	return err
}

func (r *ElasticRepository) Search(ctx context.Context, query string, skip, take uint64) ([]domain.Message, error) {
	result, err := r.client.Search().
		Index("messages").
		Query(
			elastic.NewMultiMatchQuery(query, "body").
				Fuzziness("3").
				PrefixLength(1).
				CutoffFrequency(0.0001),
		).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var messages []domain.Message
	for _, hit := range result.Hits.Hits {
		var message domain.Message
		err = json.Unmarshal(*hit.Source, &message)
		if err != nil {
			log.Println(err)
		} else {
			messages = append(messages, message)
		}
	}

	return messages, nil
}
