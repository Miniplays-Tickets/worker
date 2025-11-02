package button

import (
	"github.com/Miniplays-Tickets/worker"
	"github.com/TicketsBot-cloud/gdl/objects/interaction"
)

type Response interface {
	Type() ResponseType
	Build() interface{} // Returns the interaction response struct
	HandleDeferred(interactionData interaction.InteractionMetadata, worker *worker.Context) error
}

type ResponseType uint8

const (
	ResponseTypeMessage ResponseType = iota
	ResponseTypeEdit
	ResponseTypeModal
	ResponseTypeAck
)
