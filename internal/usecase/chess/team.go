package chess

import (
	"fmt"

	"github.com/skvdmt/chess-game-back/internal/entities"
	"github.com/skvdmt/chess-game-back/internal/entities/dto"
	"github.com/skvdmt/chess-game-back/internal/model"
)

// Team Команда.
type Team struct {
	name        string
	chessPieces *ChessPieces
	enemy       *Team
	history     *History
}

// NewTeam Конструктор.
func NewTeam(name string, history *History) *Team {
	return &Team{
		name:        name,
		chessPieces: NewChessPieces(name),
		history:     history,
	}
}

// Name Имя.
func (t *Team) Name() string {
	return t.name
}

// ChessPieces Шахматные фигуры.
func (t *Team) ChessPieces() *ChessPieces {
	return t.chessPieces
}

// Enemy Соперник.
func (t *Team) Enemy() *Team {
	return t.enemy
}

// Export Экспорт.
func (t *Team) Export() *dto.Team {
	return &dto.Team{
		Name:               t.Name(),
		OnBoardChessPieces: t.ChessPieces().ExportOnBoard(),
		EatenChessPieces:   t.ChessPieces().ExportEaten(),
	}
}

// Win Команда выиграла.
func (t *Team) Win() error {
	switch t.Name() {
	case model.White:
		return model.ErrWhiteWin
	case model.Black:
		return model.ErrBlackWin
	default:
		panic(fmt.Sprintf("unknown team name %s", t.Name()))
	}
}

// CalculatePossibleMoves Расчитать и установить возможные ходы
// каждой шахматной фигуры этой команды на доске.
func (t *Team) CalculatePossibleMoves() {
	for _, c := range *t.ChessPieces() {
		if !c.OnBoard() {
			continue
		}
		if c.Name() == model.PAWN {
			c.SetPossibleMoves(t.possibleMovesPawn(c))
			continue
		}
		var pm []*entities.Position
		if c.Name() == model.KING {
			pm = append(pm, t.possibleCastling(c)...)
		}
		for _, p := range t.beatenPositions(c) {
			if t.ChessPieces().FoundByPosition(p) != nil {
				continue
			}
			if t.checkAfterMove(c, p) {
				continue
			}
			pm = append(pm, p)
		}
		c.SetPossibleMoves(pm)
	}
}

// CanMove Может сделать ход.
func (t *Team) CanMove() bool {
	for _, c := range *t.ChessPieces() {
		if !c.OnBoard() {
			continue
		}
		if len(c.PossibleMoves()) > 0 {
			return true
		}
	}
	return false
}

// possibleMovesPawn Возможные ходы пешки.
func (t *Team) possibleMovesPawn(c ChessPiece) []*entities.Position {
	var pm []*entities.Position
	func() {
		var y1, y2 uint8
		switch t.Name() {
		case model.White:
			y1 = c.Y() + 1
			y2 = c.Y() + 2
		case model.Black:
			y1 = c.Y() - 1
			y2 = c.Y() - 2
		default:
			panic(fmt.Sprintf("unknown team name %s", t.Name()))
		}
		p := entities.NewPosition(c.X(), y1)
		if !p.OnBoard() ||
			t.ChessPieces().FoundByPosition(p) != nil ||
			t.Enemy().ChessPieces().FoundByPosition(p) != nil {
			return
		}
		if !t.checkAfterMove(c, p) {
			pm = append(pm, p)
		}
		p = entities.NewPosition(c.X(), y2)
		if c.AlreadyMove() ||
			!p.OnBoard() ||
			t.ChessPieces().FoundByPosition(p) != nil ||
			t.Enemy().ChessPieces().FoundByPosition(p) != nil {
			return
		}
		if !t.checkAfterMove(c, p) {
			pm = append(pm, p)
		}
	}()
	for _, p := range t.beatenPositions(c) {
		if t.Enemy().ChessPieces().FoundByPosition(p) == nil &&
			!t.history.PassantCapture(t, c, p) {
			continue
		}
		if t.checkAfterMove(c, p) {
			continue
		}
		pm = append(pm, p)
	}
	return pm
}

// possibleCastling Возможные рокировки.
func (t *Team) possibleCastling(c ChessPiece) []*entities.Position {
	if c.AlreadyMove() {
		return nil
	}
	if t.check() {
		return nil
	}
	var cm []*entities.Position
	r := t.ChessPieces().FoundByPosition(entities.NewPosition(1, c.Y()))
	if r != nil &&
		r.Name() == model.ROOK &&
		!r.AlreadyMove() &&
		!t.ChessPieces().ExistsByPosition(entities.NewPosition(2, c.Y())) &&
		!t.Enemy().ChessPieces().ExistsByPosition(entities.NewPosition(2, c.Y())) &&
		!t.ChessPieces().ExistsByPosition(entities.NewPosition(3, c.Y())) &&
		!t.Enemy().ChessPieces().ExistsByPosition(entities.NewPosition(3, c.Y())) &&
		!t.beaten(entities.NewPosition(3, c.Y())) &&
		!t.ChessPieces().ExistsByPosition(entities.NewPosition(4, c.Y())) &&
		!t.Enemy().ChessPieces().ExistsByPosition(entities.NewPosition(4, c.Y())) &&
		!t.beaten(entities.NewPosition(4, c.Y())) {
		cm = append(cm, entities.NewPosition(3, c.Y()))
	}
	r = t.ChessPieces().FoundByPosition(entities.NewPosition(8, c.Y()))
	if r != nil &&
		r.Name() == model.ROOK &&
		!r.AlreadyMove() &&
		!t.ChessPieces().ExistsByPosition(entities.NewPosition(6, c.Y())) &&
		!t.Enemy().ChessPieces().ExistsByPosition(entities.NewPosition(6, c.Y())) &&
		!t.beaten(entities.NewPosition(6, c.Y())) &&
		!t.ChessPieces().ExistsByPosition(entities.NewPosition(7, c.Y())) &&
		!t.Enemy().ChessPieces().ExistsByPosition(entities.NewPosition(7, c.Y())) &&
		!t.beaten(entities.NewPosition(7, c.Y())) {
		cm = append(cm, entities.NewPosition(7, c.Y()))
	}
	return cm
}

// beatenPositions Битые поля шахматной фигуры.
func (t *Team) beatenPositions(c ChessPiece) []*entities.Position {
	if c.Name() == model.PAWN {
		return t.beatenPawn(c)
	}
	if c.Name() == model.KNIGHT {
		return t.beatenKnight(c)
	}
	var bf []*entities.Position
	directions := c.Directions()
	maxRemote := c.MaxRemote()
	l := make(map[entities.Direction]struct{})
	for r := uint8(1); r <= maxRemote; r++ {
		for _, d := range directions {
			// Пропустить заблокированое направление.
			if _, o := l[d]; o {
				continue
			}
			p := t.offset(c, d, r)
			// Если позиция не на доске, заблокировать направление и продолжить итерацию.
			if !p.OnBoard() {
				l[d] = struct{}{}
				continue
			}
			// Если найдена любая фигура, заблокировать направление, но добавить битое поле.
			if t.ChessPieces().ExistsByPosition(p) ||
				t.Enemy().ChessPieces().ExistsByPosition(p) {
				l[d] = struct{}{}
			}
			bf = append(bf, p)
		}
	}
	return bf
}

// beatenPawn Битые поля пешки.
func (t *Team) beatenPawn(c ChessPiece) []*entities.Position {
	var ms []*entities.Position
	var y uint8
	switch t.Name() {
	case model.White:
		y = c.Y() + 1
	case model.Black:
		y = c.Y() - 1
	default:
		panic(fmt.Sprintf("unknown team %s", t.Name()))
	}
	ms = append(ms, entities.NewPosition(c.X()-1, y))
	ms = append(ms, entities.NewPosition(c.X()+1, y))
	var bf []*entities.Position
	for _, m := range ms {
		if m.OnBoard() {
			bf = append(bf, m)
		}
	}
	return bf
}

// beatenKnight Битые поля коня.
func (t *Team) beatenKnight(c ChessPiece) []*entities.Position {
	var ms []*entities.Position
	ms = append(ms, entities.NewPosition(c.X()+1, c.Y()+2))
	ms = append(ms, entities.NewPosition(c.X()+2, c.Y()+1))
	ms = append(ms, entities.NewPosition(c.X()-1, c.Y()-2))
	ms = append(ms, entities.NewPosition(c.X()-2, c.Y()-1))
	ms = append(ms, entities.NewPosition(c.X()+1, c.Y()-2))
	ms = append(ms, entities.NewPosition(c.X()+2, c.Y()-1))
	ms = append(ms, entities.NewPosition(c.X()-1, c.Y()+2))
	ms = append(ms, entities.NewPosition(c.X()-2, c.Y()+1))
	var bf []*entities.Position
	for _, m := range ms {
		if m.OnBoard() {
			bf = append(bf, m)
		}
	}
	return bf
}

// Offset Смещенная по направлению и дистанции позиция.
func (t *Team) offset(c ChessPiece, d entities.Direction, remote uint8) *entities.Position {
	switch d {
	case entities.Top:
		return entities.NewPosition(c.X(), c.Y()+remote)
	case entities.TopRight:
		return entities.NewPosition(c.X()+remote, c.Y()+remote)
	case entities.Right:
		return entities.NewPosition(c.X()+remote, c.Y())
	case entities.RightBottom:
		return entities.NewPosition(c.X()+remote, c.Y()-remote)
	case entities.Bottom:
		return entities.NewPosition(c.X(), c.Y()-remote)
	case entities.BottomLeft:
		return entities.NewPosition(c.X()-remote, c.Y()-remote)
	case entities.Left:
		return entities.NewPosition(c.X()-remote, c.Y())
	case entities.LeftTop:
		return entities.NewPosition(c.X()-remote, c.Y()+remote)
	}
	panic(fmt.Sprintf("unknown direction %d", d))
}

// checkAfterMove Шах после хода.
func (t *Team) checkAfterMove(c ChessPiece, to *entities.Position) bool {
	m := NewMove(t, c, to, false, t.history)
	m.Do()
	defer m.Undo()
	if t.check() {
		return true
	}
	return false
}

// beaten Битое поле.
func (t *Team) beaten(p *entities.Position) bool {
	for _, e := range *t.Enemy().ChessPieces() {
		if !e.OnBoard() {
			continue
		}
		for _, b := range t.Enemy().beatenPositions(e) {
			if b.X() == p.X() && b.Y() == p.Y() {
				return true
			}
		}
	}
	return false
}

// Check Королю поставлен шах.
func (t *Team) check() bool {
	k := t.ChessPieces().King()
	if t.beaten(k.Pos()) {
		return true
	}
	return false
}
