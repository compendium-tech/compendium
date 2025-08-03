package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/compendium-tech/compendium/college-service/internal/model"
	common "github.com/compendium-tech/compendium/common/pkg"
	"github.com/elastic/go-elasticsearch/v9"
)

type elasticSearchCollegeRepository struct {
	client *elasticsearch.Client
}

func NewElasticsearchCollegeRepository(cfg elasticsearch.Config, index string) (CollegeRepository, error) {
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %v", err)
	}

	return &elasticSearchCollegeRepository{
		client: client,
	}, nil
}

func (a *elasticSearchCollegeRepository) SearchColleges(ctx context.Context, queryText, stateOrCountry string) ([]model.College, error) {
	var query common.H

	// Build a dynamic query using a `bool` query with a `must` clause.
	// This allows combining multiple conditions.
	boolQuery := make(map[string]any)
	var must []map[string]any

	if queryText != "" {
		must = append(must, common.H{
			"semantic": common.H{
				"field": "description",
				"query": queryText,
			},
		})
	}

	if stateOrCountry != "" {
		must = append(must, common.H{
			"match": common.H{
				"state_or_country": strings.ToLower(stateOrCountry),
			},
		})
	}

	// If both search parameters are empty, perform a match_all query.
	if len(must) == 0 {
		query = common.H{
			"match_all": common.H{},
		}
	} else {
		// If either parameter is present, use the bool query.
		boolQuery["must"] = must
		query = common.H{
			"bool": boolQuery,
		}
	}

	searchBody, _ := json.Marshal(common.H{
		"query": query,
	})

	searchRes, err := a.client.Search(
		a.client.Search.WithContext(ctx),
		a.client.Search.WithIndex("colleges"),
		a.client.Search.WithBody(bytes.NewReader(searchBody)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform search: %w", err)
	}
	defer searchRes.Body.Close()

	if searchRes.IsError() {
		return nil, fmt.Errorf("error searching: %s", searchRes.String())
	}

	var searchResult map[string]any
	if err := json.NewDecoder(searchRes.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	var colleges []model.College
	// Safely extract hits and iterate through them to populate the colleges slice.
	if hits, ok := searchResult["hits"].(common.H)["hits"].([]any); ok {
		for _, hit := range hits {
			source, ok := hit.(common.H)["_source"].(common.H)
			if !ok {
				continue
			}
			college := model.College{
				Name:           source["name"].(string),
				City:           source["city"].(string),
				StateOrCountry: source["state_or_country"].(string),
				Description:    source["description"].(string),
			}
			colleges = append(colleges, college)
		}
	}

	return colleges, nil
}
