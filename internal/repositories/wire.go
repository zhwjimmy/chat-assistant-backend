package repositories

import (
	"github.com/google/wire"
)

// RepositorySet provides all repository dependencies
var RepositorySet = wire.NewSet(
	NewUserRepository,
	NewConversationRepository,
	NewMessageRepository,
	NewTagRepository,
	NewElasticsearchRepository,
)
