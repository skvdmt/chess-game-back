package model

import "errors"

// Errors Глобальный канал ошибок.
var Errors chan error

var (
	ErrWaitBothPlayers        = errors.New("wait both players")
	ErrWaitWhitePlayer        = errors.New("wait white player")
	ErrWaitBlackPlayer        = errors.New("wait black player")
	ErrGameOver               = errors.New("game over")
	ErrTimeOver               = errors.New("time over")
	ErrWhiteWin               = errors.New("white win")
	ErrBlackWin               = errors.New("black win")
	ErrDraw                   = errors.New("draw")
	ErrGameNotOver            = errors.New("game not over")
	ErrNotYourMove            = errors.New("it's not your move now")
	ErrYouAreSpectator        = errors.New("you are spectator")
	ErrChessPieceNotFound     = errors.New("chess Piece Not Found")
	ErrMoveImpossible         = errors.New("move impossible")
	ErrAttemptsOfferADrawOver = errors.New("attempts offer a draw over")
	ErrOfferADrawAlreadySent  = errors.New("offer a draw already sent")
	ErrNoOfferADraw           = errors.New("did not offer a draw")

	ErrRulesRook   = errors.New("rook can't move like that")
	ErrRulesKnight = errors.New("knight can't move like that")
	ErrRulesBishop = errors.New("bishop can't move like that")
	ErrRulesQueen  = errors.New("queen can't move like that")
	ErrRulesKing   = errors.New("king can't move like that")
	ErrRulesPawn   = errors.New("pawn can't move like that")
)
