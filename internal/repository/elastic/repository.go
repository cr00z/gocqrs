package elastic

import (
	"context"
	"encoding/json"
	"github.com/cr00z/gocqrs/internal/domain"
	"github.com/olivere/elastic"
	"log"
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

func (r *ElasticRepository) Insert(ctx context.Context, meow domain.Meow) error {
	_, err := r.client.Index().
		Index("meows").
		Type("meow").
		Id(meow.ID).
		BodyJson(meow).
		Refresh("wait_for").Do(ctx)
	return err
}

func (r *ElasticRepository) Search(ctx context.Context, query string, skip, take uint64) ([]domain.Meow, error) {
	result, err := r.client.Search().
		Index("meows").
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

	var meows []domain.Meow
	for _, hit := range result.Hits.Hits {
		var meow domain.Meow
		err = json.Unmarshal(*hit.Source, &meow)
		if err != nil {
			log.Println(err)
		} else {
			meows = append(meows, meow)
		}
	}

	return meows, nil
}
