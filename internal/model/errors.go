package model

import "errors"

// Errors Глобальный канал ошибок.
var Errors chan error

var (
	ErrWaitBothPlayers         = errors.New("Wait Both Players")
	ErrWaitWhitePlayer         = errors.New("Wait White Player")
	ErrWaitBlackPlayer         = errors.New("Wait Black Player")
	ErrGameOver                = errors.New("Game Over")
	ErrTimeOver                = errors.New("time over")
	ErrWhiteWin                = errors.New("White Win")
	ErrBlackWin                = errors.New("Black Win")
	ErrDraw                    = errors.New("Draw")
	ErrGameNotOver             = errors.New("Game not over")
	ErrNotYourMove             = errors.New("It's not your move now")
	ErrYouAreSpectator         = errors.New("You are spectator")
	ErrChessPieceNotFound      = errors.New("Chess Piece Not Found")
	ErrMoveImpossible          = errors.New("Move Impossible")
	ErrAttemptsOfferADrawOver  = errors.New("Attempts Offer A Draw Over")
	ErrOfferADrawAlreadySended = errors.New("Offer A Draw Already Sended")
	ErrNoOfferADraw            = errors.New("Did not offer a draw")

	ErrRulesRook   = errors.New("rook can't move like that")
	ErrRulesKnight = errors.New("knight can't move like that")
	ErrRulesBishop = errors.New("bishop can't move like that")
	ErrRulesQueen  = errors.New("queen can't move like that")
	ErrRulesKing   = errors.New("king can't move like that")
	ErrRulesPawn   = errors.New("pawn can't move like that")
)
