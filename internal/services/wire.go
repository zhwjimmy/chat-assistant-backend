package services

import (
	"github.com/google/wire"
)

// ServiceSet provides all service dependencies
var ServiceSet = wire.NewSet(
	NewUserService,
	NewConversationService,
	NewMessageService,
	NewTagService,
	NewSearchService,
	NewSyncService,
)
