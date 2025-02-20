package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crosschain/testutil"
	"github.com/cosmos/cosmos-sdk/x/crosschain/types"
)

type TestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (s *TestSuite) SetupTest() {
	app := simapp.Setup(s.T(), false, true)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})

	app.CrossChainKeeper.SetParams(ctx, types.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.CrossChainKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	s.app = app
	s.ctx = ctx
	s.queryClient = queryClient
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestIncrSendSequence() {
	beforeSequence := s.app.CrossChainKeeper.GetSendSequence(s.ctx, sdk.ChannelID(1))

	s.app.CrossChainKeeper.IncrSendSequence(s.ctx, sdk.ChannelID(1))

	afterSequence := s.app.CrossChainKeeper.GetSendSequence(s.ctx, sdk.ChannelID(1))

	s.Require().EqualValues(afterSequence, beforeSequence+1)
}

func (s *TestSuite) TestIncrReceiveSequence() {
	beforeSequence := s.app.CrossChainKeeper.GetReceiveSequence(s.ctx, sdk.ChannelID(1))

	s.app.CrossChainKeeper.IncrReceiveSequence(s.ctx, sdk.ChannelID(1))

	afterSequence := s.app.CrossChainKeeper.GetReceiveSequence(s.ctx, sdk.ChannelID(1))

	s.Require().EqualValues(afterSequence, beforeSequence+1)
}

func (s *TestSuite) TestRegisterChannel() {
	testChannelName := "test channel"
	testChannelId := sdk.ChannelID(100)

	err := s.app.CrossChainKeeper.RegisterChannel(testChannelName, testChannelId, &testutil.MockCrossChainApplication{})

	s.Require().NoError(err)

	app := s.app.CrossChainKeeper.GetCrossChainApp(testChannelId)
	s.Require().NotNil(app)

	// check duplicate name
	err = s.app.CrossChainKeeper.RegisterChannel(testChannelName, testChannelId, app)
	s.Require().ErrorContains(err, "duplicated channel name")

	// check duplicate channel id
	err = s.app.CrossChainKeeper.RegisterChannel("another channel", testChannelId, app)
	s.Require().ErrorContains(err, "duplicated channel id")

	// check nil app
	err = s.app.CrossChainKeeper.RegisterChannel("another channel", sdk.ChannelID(101), nil)
	s.Require().ErrorContains(err, "nil cross chain app")
}

func (s *TestSuite) TestSetChannelSendPermission() {
	s.app.CrossChainKeeper.SetChannelSendPermission(s.ctx, sdk.ChainID(1), sdk.ChannelID(1), sdk.ChannelAllow)

	permission := s.app.CrossChainKeeper.GetChannelSendPermission(s.ctx, sdk.ChainID(1), sdk.ChannelID(1))
	s.Require().EqualValues(sdk.ChannelAllow, permission)
}
