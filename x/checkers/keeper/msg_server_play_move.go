package keeper

import (
	"context"
	"strconv"

	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Fetch the stored game information
	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
	}

	// Check that the game has not finished yet
	if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER] {
		return nil, types.ErrGameFinished
	}

	// Is the player Legitimate?
	isBlack := storedGame.Black == msg.Creator
	isRed := storedGame.Red == msg.Creator
	var player rules.Player
	if !isBlack && !isRed {
		return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
	} else if isBlack && isRed {
		player = rules.StringPieces[storedGame.Turn].Player
	} else if isBlack {
		player = rules.BLACK_PLAYER
	} else {
		player = rules.RED_PLAYER
	}

	// Instantiate the board in order to implement the rules:
	game, err := storedGame.ParseGame()
	if err != nil {
		panic(err.Error())
	}

	// Is it the player's turn?
	if !game.TurnIs(player) {
		return nil, sdkerrors.Wrapf(types.ErrNotPlayerTurn, "%s", player)
	}

	// Properly conduct the move
	captured, moveErr := game.Move(
		rules.Pos{
			X: int(msg.FromX),
			Y: int(msg.FromY),
		},
		rules.Pos{
			X: int(msg.ToX),
			Y: int(msg.ToY),
		},
	)
	if moveErr != nil {
		return nil, sdkerrors.Wrapf(types.ErrWrongMove, moveErr.Error())
	}

	// Prepare the updated board to be stored and store the information
	storedGame.Winner = rules.PieceStrings[game.Winner()]

	lastBoard := game.String()
	if storedGame.Winner == rules.PieceStrings[rules.NO_PLAYER] {
		storedGame.Board = lastBoard
	} else {
		storedGame.Board = ""
	}

	storedGame.Turn = rules.PieceStrings[game.Turn]
	k.Keeper.SetStoredGame(ctx, storedGame)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.GameCreatedEventType,
			sdk.NewAttribute(types.MovePlayedEventType, msg.Creator),
			sdk.NewAttribute(types.MovePlayedEventGameIndex, msg.GameIndex),
			sdk.NewAttribute(types.MovePlayedEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
			sdk.NewAttribute(types.MovePlayedEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
			sdk.NewAttribute(types.MovePlayedEventWinner, rules.PieceStrings[game.Winner()]),
			sdk.NewAttribute(types.MovePlayedEventBoard, lastBoard),
		),
	)

	// Return relevant information regarding the move's result
	return &types.MsgPlayMoveResponse{
		CapturedX: int32(captured.X),
		CapturedY: int32(captured.Y),
		Winner:    rules.PieceStrings[game.Winner()],
	}, nil
}
