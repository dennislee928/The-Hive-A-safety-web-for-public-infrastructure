package cap

import (
	"github.com/erh-safety-system/poc/internal/route1"
)

// CAPMessageAdapter adapts CAPMessage to route1.CAPMessageInterface
type CAPMessageAdapter struct {
	msg *CAPMessage
}

// NewCAPMessageAdapter creates a new adapter
func NewCAPMessageAdapter(msg *CAPMessage) route1.CAPMessageInterface {
	return &CAPMessageAdapter{msg: msg}
}

// GetIdentifier returns the identifier
func (a *CAPMessageAdapter) GetIdentifier() string {
	return a.msg.Identifier
}

// GetSender returns the sender
func (a *CAPMessageAdapter) GetSender() string {
	return a.msg.Sender
}

// GetSent returns the sent time
func (a *CAPMessageAdapter) GetSent() string {
	return a.msg.Sent
}

// GetStatus returns the status
func (a *CAPMessageAdapter) GetStatus() string {
	return a.msg.Status
}

// GetMsgType returns the message type
func (a *CAPMessageAdapter) GetMsgType() string {
	return a.msg.MsgType
}

// GetScope returns the scope
func (a *CAPMessageAdapter) GetScope() string {
	return a.msg.Scope
}

// GetInfoBlocks returns the info blocks
func (a *CAPMessageAdapter) GetInfoBlocks() []route1.InfoBlock {
	blocks := make([]route1.InfoBlock, len(a.msg.Info))
	for i, info := range a.msg.Info {
		blocks[i] = route1.InfoBlock{
			Language:    info.Language,
			Headline:    info.Headline,
			Description: info.Description,
			Instruction: info.Instruction,
			Expires:     info.Expires,
		}
	}
	return blocks
}

// GetArea returns the area info
func (a *CAPMessageAdapter) GetArea() route1.AreaInfo {
	return route1.AreaInfo{
		ZoneID:   a.msg.Area.ZoneID,
		ZoneType: a.msg.Area.ZoneType,
	}
}

// GetSignature returns the signature info
func (a *CAPMessageAdapter) GetSignature() *route1.SignatureInfo {
	if a.msg.Signature == nil {
		return nil
	}
	return &route1.SignatureInfo{
		Algorithm: a.msg.Signature.Algorithm,
		Value:     a.msg.Signature.Value,
	}
}

