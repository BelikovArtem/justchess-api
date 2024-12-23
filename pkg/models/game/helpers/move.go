package helpers

import (
	"justchess/pkg/models/game/enums"
	"time"
)

// Move is used to store completed moves in a database.
type Move struct {
	To               Pos             `json:"to"`
	From             Pos             `json:"from"`
	MoveType         enums.MoveType  `json:"-"`
	IsCheck          bool            `json:"-"`
	IsCapture        bool            `json:"-"`
	IsCheckmate      bool            `json:"-"`
	TimeLeft         time.Duration   `json:"-"`
	PromotionPayload enums.PieceType `json:"pp"`
}

// ToLAN serializes the given move into Long Algebraic Notation.
func (m Move) ToLAN(pt enums.PieceType) string {
	lan := ""

	switch m.MoveType {
	case enums.ShortCastling:
		lan = "0-0"

	case enums.LongCastling:
		lan = "0-0-0"

	case enums.Promotion:
		if m.IsCapture {
			lan += m.From.String() + "x" + m.To.String()
		} else {
			lan += m.From.String() + "-" + m.To.String()
		}
		lan += "=" + m.PromotionPayload.String()

	default:
		if pt != enums.Pawn {
			lan = pt.String()
		}
		if m.IsCapture {
			lan += m.From.String() + "x" + m.To.String()
		} else {
			lan += m.From.String() + "-" + m.To.String()
		}
	}
	if m.IsCheck {
		if m.IsCheckmate {
			lan += "#"
		} else {
			lan += "+"
		}
	}
	return lan
}

// MoveDTO is stored in a database and displayed on a frontend.
type MoveDTO struct {
	UCI        string            `json:"uci"`
	LAN        string            `json:"lan"`
	FEN        string            `json:"fen"`
	TimeLeft   time.Duration     `json:"timeLeft"`
	ValidMoves map[string]string `json:"vm"`
}

// PossibleMove represents player`s possible moves.
type PossibleMove struct {
	To       Pos
	MoveType enums.MoveType
}

// NewPM creates a new PossibleMove.
func NewPM(to Pos, mt enums.MoveType) PossibleMove {
	return PossibleMove{
		To:       to,
		MoveType: mt,
	}
}
