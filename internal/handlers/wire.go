package handlers

import (
	"github.com/google/wire"
)

// HandlerSet provides all handler dependencies
var HandlerSet = wire.NewSet(
	NewUserHandler,
	NewConversationHandler,
	NewMessageHandler,
	NewTagHandler,
	NewSearchHandler,
)
