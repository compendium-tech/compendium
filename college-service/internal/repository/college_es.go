package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/compendium-tech/compendium/college-service/internal/model"
	common "github.com/compendium-tech/compendium/common/pkg"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/ztrue/tracerr"
)

type elasticSearchCollegeRepository struct {
	client *elasticsearch.Client
}

func NewElasticsearchCollegeRepository(client *elasticsearch.Client) CollegeRepository {
	return &elasticSearchCollegeRepository{
		client: client,
	}
}

func (a *elasticSearchCollegeRepository) SearchColleges(
	ctx context.Context, semanticSearchText,
	stateOrCountry string, pageIndex, pageSize int) ([]model.College, error) {
	// Set default pagination values if they are invalid.
	if pageIndex < 0 {
		pageIndex = 0
	}
	if pageSize < 1 {
		pageSize = 10 // A reasonable default page size
	}

	from := pageIndex * pageSize

	var query common.H

	// Build a dynamic query using a `bool` query with a `must` clause.
	// This allows combining multiple conditions.
	boolQuery := make(map[string]any)
	var must []map[string]any

	if semanticSearchText != "" {
		must = append(must, common.H{
			"semantic": common.H{
				"field": "description",
				"query": semanticSearchText,
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
		"from":  from,
		"size":  pageSize,
	})

	searchRes, err := a.client.Search(
		a.client.Search.WithContext(ctx),
		a.client.Search.WithIndex("colleges"),
		a.client.Search.WithBody(bytes.NewReader(searchBody)),
	)
	if err != nil {
		return nil, tracerr.Errorf("failed to perform search: %v", err)
	}
	defer searchRes.Body.Close()

	if searchRes.IsError() {
		return nil, tracerr.Errorf("error searching: %s", searchRes.String())
	}

	var searchResult map[string]any
	if err := json.NewDecoder(searchRes.Body).Decode(&searchResult); err != nil {
		return nil, tracerr.Errorf("failed to decode search response: %v", err)
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
