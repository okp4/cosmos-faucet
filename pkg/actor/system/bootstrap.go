package system

import (
	"okp4/cosmos-faucet/pkg/cosmos"
	"okp4/cosmos-faucet/pkg/faucet"

	"github.com/asynkron/protoactor-go/actor"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/credentials"
)

func BootstrapActors(
	chainID string,
	privKey crypto.PrivKey,
	sendAmount types.Coins,
	grpcAddress string,
	tls credentials.TransportCredentials,
) (*actor.RootContext, *actor.PID) {
	cosmosClientProps := actor.PropsFromProducer(func() actor.Actor {
		grpcClient, err := cosmos.NewGrpcClient(grpcAddress, tls)
		if err != nil {
			log.Panic().Err(err).Msg("❌ Could not create grpc client")
		}

		return grpcClient
	})

	txHandlerProps := actor.PropsFromProducer(func() actor.Actor {
		return cosmos.NewTxHandler(
			cosmos.WithChainID(chainID),
			cosmos.WithPrivateKey(privKey),
			cosmos.WithTxConfig(simapp.MakeTestEncodingConfig().TxConfig),
			cosmos.WithCosmosClientProps(cosmosClientProps),
		)
	})

	actorCTX := actor.NewActorSystem().Root
	return actorCTX, actorCTX.Spawn(actor.PropsFromProducer(func() actor.Actor {
		return faucet.NewFaucet(
			faucet.WithAmount(sendAmount),
			faucet.WithAddress(types.AccAddress(privKey.PubKey().Address())),
			faucet.WithTxHandlerProps(txHandlerProps),
		)
	}))
}
