package keeper_test

import (
	"testing"

	"github.com/alice/checkers/x/checkers/testutil"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestPlayMoveUpToWinner(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)

	testutil.PlayAllMoves(t, msgServer, context, "1", bob, carol, testutil.Game1Moves)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: types.NoFifoIndex,
		FifoTailIndex: types.NoFifoIndex,
	}, systemInfo)

	game, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		Winner:      "b",
		Deadline:    types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		MoveCount:   40,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  types.NoFifoIndex,
		Wager:       45,
	}, game)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.Equal(t, event.Type, "move-played")
	require.EqualValues(t, []sdk.Attribute{
		{Key: "creator", Value: bob},
		{Key: "game-index", Value: "1"},
		{Key: "captured-x", Value: "2"},
		{Key: "captured-y", Value: "5"},
		{Key: "winner", Value: "b"},
		{Key: "board", Value: "*b*b****|**b*b***|*****b**|********|***B****|********|*****b**|********"},
	}, event.Attributes[(len(testutil.Game1Moves)-1)*6:])
}

func TestPlayMoveUpToWinnerCalledBank(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMoveWithMock(t)
	defer ctrl.Finish()
	payBob := escrow.ExpectPay(context, bob, 45).Times(1)
	payCarol := escrow.ExpectPay(context, carol, 45).Times(1).After(payBob)
	escrow.ExpectRefund(context, bob, 90).Times(1).After(payCarol)

	testutil.PlayAllMoves(t, msgServer, context, "1", bob, carol, testutil.Game1Moves)
}
