package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	elastic "github.com/olivere/elastic/v7"
)

var (
	errNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	DocType     string  `json:"doc_type"`
}

func NewElasticRepository(url string) (*elasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ERROR ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC INFO ", log.LstdFlags)),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {

}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().
		Index("catalog").
		Id(p.ID).
		BodyJson(productDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			DocType:     "product",
		}).
		Do(ctx)

	return err
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("catalog").
		Id(id).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	if !res.Found {
		return nil, errNotFound
	}

	p := productDocument{}
	if err = json.Unmarshal(res.Source, &p); err != nil {
		return nil, err
	}
	return &Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, err
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).Size(int(take)).
		Do(ctx)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err = json.Unmarshal(hit.Source, &p); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, err
}

func (r *elasticRepository) ListProductWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	items := []*elastic.MultiGetItem{}

	for _, id := range ids {
		items = append(
			items,
			elastic.NewMultiGetItem().
				Index("catalog").
				Id(id),
		)
	}

	res, err := r.client.MultiGet().
		Add(items...).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []Product{}
	for _, doc := range res.Docs {
		p := productDocument{}
		if err = json.Unmarshal(doc.Source, &p); err == nil {
			products = append(products, Product{
				ID:          doc.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}

	return products, err
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err = json.Unmarshal(hit.Source, &p); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}

	}
	return products, err
}
