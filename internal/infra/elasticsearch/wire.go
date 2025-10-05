package elasticsearch

import (
	"github.com/google/wire"
)

// ElasticsearchSet provides all Elasticsearch-related dependencies
var ElasticsearchSet = wire.NewSet(
	NewElasticsearchClientFromConfig,
	NewElasticsearchIndexerFromClient,
	NewElasticsearchClient,
	NewElasticsearchIndexName,
)
