package cosmos

import (
	"fmt"
	"okp4/cosmos-faucet/pkg/actor/message"
	"okp4/cosmos-faucet/test/mock"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	signing2 "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/stretchr/testify/mock"
)

var (
	mnemonic = "orient few plate lawsuit pumpkin shallow wild relax mosquito ill stool pitch orbit develop ginger bachelor slot useful worry arena bamboo combine arch envelope"
	chainID  = "testnet-okp4-1"
	txConfig = simapp.MakeTestEncodingConfig().TxConfig

	msgs = []types.Msg{
		banktypes.NewMsgSend(
			types.AccAddress("from"),
			types.AccAddress("to"),
			types.NewCoins(types.NewInt64Coin("uknow", 1000000)),
		),
		banktypes.NewMsgSend(
			types.AccAddress("from"),
			types.AccAddress("to"),
			types.NewCoins(types.NewInt64Coin("uknow", 5000000)),
		),
	}
	memo       = "Sent from tests"
	gasLimit   = uint64(200000)
	feeAmount  = types.NewCoins(types.NewInt64Coin("uknow", 50000))
	signerData = signing.SignerData{
		ChainID:       chainID,
		AccountNumber: 10,
		Sequence:      54,
	}
	signData = &signing2.SingleSignatureData{
		SignMode:  txConfig.SignModeHandler().DefaultMode(),
		Signature: []uint8{119, 228, 243, 201, 202, 94, 114, 31, 59, 53, 109, 198, 174, 101, 91, 211, 189, 4, 30, 39, 142, 187, 123, 199, 255, 129, 19, 91, 228, 229, 123, 202, 9, 154, 140, 241, 253, 186, 250, 153, 62, 178, 183, 126, 241, 8, 71, 161, 53, 153, 62, 202, 241, 22, 67, 132, 37, 78, 240, 41, 83, 157, 161, 16},
	}
	txRaw = []uint8{10, 203, 1, 10, 91, 10, 28, 47, 99, 111, 115, 109, 111, 115, 46, 98, 97, 110, 107, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 83, 101, 110, 100, 18, 59, 10, 20, 99, 111, 115, 109, 111, 115, 49, 118, 101, 101, 120, 55, 109, 103, 122, 116, 56, 51, 99, 117, 18, 17, 99, 111, 115, 109, 111, 115, 49, 119, 51, 104, 115, 106, 116, 116, 114, 102, 113, 26, 16, 10, 5, 117, 107, 110, 111, 119, 18, 7, 49, 48, 48, 48, 48, 48, 48, 10, 91, 10, 28, 47, 99, 111, 115, 109, 111, 115, 46, 98, 97, 110, 107, 46, 118, 49, 98, 101, 116, 97, 49, 46, 77, 115, 103, 83, 101, 110, 100, 18, 59, 10, 20, 99, 111, 115, 109, 111, 115, 49, 118, 101, 101, 120, 55, 109, 103, 122, 116, 56, 51, 99, 117, 18, 17, 99, 111, 115, 109, 111, 115, 49, 119, 51, 104, 115, 106, 116, 116, 114, 102, 113, 26, 16, 10, 5, 117, 107, 110, 111, 119, 18, 7, 53, 48, 48, 48, 48, 48, 48, 18, 15, 83, 101, 110, 116, 32, 102, 114, 111, 109, 32, 116, 101, 115, 116, 115, 18, 104, 10, 80, 10, 70, 10, 31, 47, 99, 111, 115, 109, 111, 115, 46, 99, 114, 121, 112, 116, 111, 46, 115, 101, 99, 112, 50, 53, 54, 107, 49, 46, 80, 117, 98, 75, 101, 121, 18, 35, 10, 33, 2, 207, 77, 102, 22, 143, 32, 134, 42, 254, 197, 154, 121, 37, 128, 189, 112, 30, 242, 85, 27, 59, 123, 56, 120, 111, 252, 57, 53, 131, 79, 175, 49, 18, 4, 10, 2, 8, 1, 24, 54, 18, 20, 10, 14, 10, 5, 117, 107, 110, 111, 119, 18, 5, 53, 48, 48, 48, 48, 16, 192, 154, 12, 26, 64, 119, 228, 243, 201, 202, 94, 114, 31, 59, 53, 109, 198, 174, 101, 91, 211, 189, 4, 30, 39, 142, 187, 123, 199, 255, 129, 19, 91, 228, 229, 123, 202, 9, 154, 140, 241, 253, 186, 250, 153, 62, 178, 183, 126, 241, 8, 71, 161, 53, 153, 62, 202, 241, 22, 67, 132, 37, 78, 240, 41, 83, 157, 161, 16}
)

func TestOptions(t *testing.T) {
	Convey("Given a set of options", t, func() {
		privKey, err := ParseMnemonic(mnemonic)
		if err != nil {
			panic(err)
		}
		opts := []Option{
			WithMnemonicMust(mnemonic),
			WithChainID(chainID),
			WithTxConfig(txConfig),
			WithCosmosClientProps(&actor.Props{}),
		}

		Convey("When creating the TxHandler with the options", func() {
			txHandler := NewTxHandler(opts...)

			Convey("Then the returned txHandler should be configured accordingly", func() {
				So(txHandler.privKey, ShouldResemble, privKey)
				So(txHandler.address, ShouldEqual, types.AccAddress(privKey.PubKey().Address()).String())
				So(txHandler.chainID, ShouldEqual, chainID)
				So(txHandler.config, ShouldResemble, txConfig)
				So(txHandler.signMode, ShouldResemble, txConfig.SignModeHandler().DefaultMode())
				So(txHandler.cosmosClientProps, ShouldResemble, &actor.Props{})
				So(txHandler.cosmosClient, ShouldBeNil)
			})
		})
	})
}

func TestBuildUnsignedTx(t *testing.T) {
	Convey("Given a TxHandler and some transaction parameters", t, func() {
		txHandler := NewTxHandler(
			WithMnemonicMust(mnemonic),
			WithChainID(chainID),
			WithTxConfig(txConfig),
		)

		Convey("When building an unsigned transaction", func() {
			unsignedTx, err := txHandler.BuildUnsignedTx(msgs, memo, gasLimit, feeAmount)

			Convey("Then the transaction should be successfully built", func() {
				So(err, ShouldBeNil)
				So(unsignedTx, ShouldNotBeNil)
			})

			Convey("And the transaction should be correctly configured", func() {
				So(unsignedTx.GetTx().GetGas(), ShouldEqual, gasLimit)
				So(unsignedTx.GetTx().GetFee().String(), ShouldEqual, feeAmount.String())
				So(unsignedTx.GetTx().GetMemo(), ShouldEqual, memo)
				So(unsignedTx.GetTx().GetMsgs(), ShouldResemble, msgs)
			})
		})
	})
}

func TestSignTx(t *testing.T) {
	Convey("Given a TxHandler and an unsigned transaction", t, func() {
		txHandler := NewTxHandler(
			WithMnemonicMust(mnemonic),
			WithChainID(chainID),
			WithTxConfig(txConfig),
		)
		unsignedTx, err := txHandler.BuildUnsignedTx(msgs, memo, gasLimit, feeAmount)
		if err != nil {
			panic(err)
		}

		Convey("When signing the transaction", func() {
			signedTx, err := txHandler.SignTx(unsignedTx, signerData)

			Convey("Then the transaction should be signed", func() {
				So(err, ShouldBeNil)

				signatures, err := signedTx.GetSignaturesV2()
				So(signatures, ShouldNotBeEmpty)
				So(err, ShouldBeNil)

				Convey("And the signature should be the same than the signer public key", func() {
					So(signatures[0].PubKey.String(), ShouldEqual, txHandler.privKey.PubKey().String())
				})

				Convey("And the account sequence should be correctly set", func() {
					So(signatures[0].Sequence, ShouldEqual, signerData.Sequence)
				})

				Convey("And the signature data shall correspond to the tx", func() {
					So(signatures[0].Data, ShouldResemble, signData)
				})
			})
		})
	})
}

func TestEncodeTx(t *testing.T) {
	Convey("Given a TxHandler and a signed transaction", t, func() {
		txHandler := NewTxHandler(
			WithMnemonicMust(mnemonic),
			WithChainID(chainID),
			WithTxConfig(txConfig),
		)
		unsignedTx, err := txHandler.BuildUnsignedTx(msgs, memo, gasLimit, feeAmount)
		if err != nil {
			panic(err)
		}
		signedTx, err := txHandler.SignTx(unsignedTx, signerData)
		if err != nil {
			panic(err)
		}

		Convey("When Encoding the transaction", func() {
			encodedTx, err := txHandler.EncodeTx(signedTx)

			Convey("Then the transaction should be properly encoded", func() {
				So(err, ShouldBeNil)
				So(encodedTx, ShouldResemble, txRaw)
			})
		})
	})
}

func TestParseMnemonic(t *testing.T) {
	Convey("Given a well formed mnemonic", t, func() {
		mnemonic := "nasty random alter chronic become keen stadium test chaos fashion during claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard"

		Convey("When generating its corresponding private key", func() {
			privKey, err := ParseMnemonic(mnemonic)

			Convey("Then The private key have been successfully decoded", func() {
				So(privKey, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a malformed mnemonic", t, func() {
		mnemonic := "nasty random alter chronic become keen stadium test chaos fashion durin claim rug thing trade swap bleak shuffle bronze gun tobacco length aim hazard"

		Convey("When generating its corresponding private key", func() {
			privKey, err := ParseMnemonic(mnemonic)

			Convey("Then the private key shall not be decoded", func() {
				So(privKey, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err, ShouldBeError, "Invalid mnemonic")
			})
		})
	})
}

func TestStarted(t *testing.T) {
	Convey("Given a tx handler actor", t, func() {
		txHandler := NewTxHandler(
			WithCosmosClientProps(&actor.Props{}),
		)

		Convey("When receiving a Started message", func() {
			mockedContext := &mock.ActorContext{}
			mockedContext.On("Message").Return(&actor.Started{})
			mockedContext.On("Spawn", Anything).Return(&actor.PID{})
			txHandler.Receive(mockedContext)

			Convey("Then the cosmos client actor should be spawned", func() {
				mockedContext.AssertCalled(t, "Message")
				mockedContext.AssertCalled(t, "Spawn", Anything)
				So(txHandler.cosmosClient, ShouldNotBeNil)
			})
		})
	})
}

func TestMakeTx(t *testing.T) {
	Convey("Given a tx handler actor", t, func() {
		txHandler := NewTxHandler(
			WithMnemonicMust(mnemonic),
			WithChainID(chainID),
			WithTxConfig(txConfig),
		)
		txHandler.cosmosClient = &actor.PID{Id: "client"}

		Convey("And a MakeTx message with an exceeded deadline", func() {
			msg := &message.MakeTx{
				Deadline: time.Now().Add(-time.Second),
			}

			Convey("When receiving the message", func() {
				mockedContext := &mock.ActorContext{}
				mockedContext.On("Message").Return(msg)
				txHandler.Receive(mockedContext)

				Convey("Then it should not make the transaction", func() {
					mockedContext.AssertCalled(t, "Message")
					mockedContext.AssertNotCalled(t, "Spawn", Anything)
					mockedContext.AssertNotCalled(t, "Send", Anything)
				})
			})
		})

		Convey("And a valid MakeTx message", func() {
			subscriber := &actor.PID{Id: "subscriber"}
			msg := &message.MakeTx{
				Deadline:     time.Now().Add(time.Second),
				TxSubscriber: subscriber,
				Msgs:         msgs,
				Memo:         memo,
				GasLimit:     gasLimit,
				FeeAmount:    feeAmount,
			}

			Convey("And a cosmos client getting error on GetAccount message", func() {
				mockedContext := &mock.ActorContext{}
				mockedContext.On("RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.GetAccount"), Anything).
					Return(
						mock.MakeFuture(
							&message.GetAccountResponse{
								Account: &auth.BaseAccount{
									AccountNumber: signerData.AccountNumber,
									Sequence:      signerData.Sequence,
								},
							},
							nil,
						),
					)

				Convey("When receiving the message", func() {
					mockedContext.On("Message").Return(msg)

					Convey("Then it should crash", func() {
						So(func() { txHandler.Receive(mockedContext) }, ShouldPanic)
						mockedContext.AssertCalled(t, "Message")
						mockedContext.AssertCalled(t, "RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.GetAccount"), Anything)
						mockedContext.AssertNotCalled(t, "RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.BroadcastTx"), Anything)
					})
				})
			})

			Convey("And a cosmos client getting error on BroadcastTx message", func() {
				mockedContext := &mock.ActorContext{}
				mockedContext.On("RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.GetAccount"), Anything).
					Return(
						mock.MakeFuture(
							&message.GetAccountResponse{
								Account: &auth.BaseAccount{
									AccountNumber: signerData.AccountNumber,
									Sequence:      signerData.Sequence,
								},
							},
							nil,
						),
					)
				mockedContext.On("RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.BroadcastTx"), Anything).
					Return(
						mock.MakeFuture(
							nil,
							fmt.Errorf("error"),
						),
					)

				Convey("When receiving the message", func() {
					mockedContext.On("Message").Return(msg)

					Convey("Then it should crash", func() {
						So(func() { txHandler.Receive(mockedContext) }, ShouldPanic)
						mockedContext.AssertCalled(t, "Message")
						mockedContext.AssertCalled(t, "RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.GetAccount"), Anything)
						mockedContext.AssertCalled(t, "RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.BroadcastTx"), Anything)
					})
				})
			})

			Convey("And a finally working cosmos client", func() {
				var broadcastMessage interface{}
				var broadcastRespMessage interface{}
				mockedContext := &mock.ActorContext{}
				mockedContext.On("RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.GetAccount"), Anything).
					Return(
						mock.MakeFuture(
							&message.GetAccountResponse{
								Account: &auth.BaseAccount{
									AccountNumber: signerData.AccountNumber,
									Sequence:      signerData.Sequence,
								},
							},
							nil,
						),
					)
				mockedContext.On("RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.BroadcastTx"), Anything).
					Run(func(args Arguments) {
						broadcastMessage = args.Get(1)
					}).
					Return(
						mock.MakeFuture(
							&message.BroadcastTxResponse{
								TxResponse: &types.TxResponse{},
							},
							nil,
						),
					)
				mockedContext.On("Send", msg.TxSubscriber, AnythingOfType("*message.BroadcastTxResponse"), Anything).
					Run(func(args Arguments) {
						broadcastRespMessage = args.Get(1)
					})

				Convey("When receiving the message", func() {
					mockedContext.On("Message").Return(msg)
					txHandler.Receive(mockedContext)

					Convey("Then it should succeed", func() {
						mockedContext.AssertCalled(t, "Message")
						mockedContext.AssertCalled(t, "RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.GetAccount"), Anything)
						mockedContext.AssertCalled(t, "RequestFuture", txHandler.cosmosClient, AnythingOfType("*message.BroadcastTx"), Anything)
						mockedContext.AssertCalled(t, "Send", msg.TxSubscriber, AnythingOfType("*message.BroadcastTxResponse"), Anything)
						So(broadcastMessage, ShouldHaveSameTypeAs, &message.BroadcastTx{})
						So(broadcastMessage.(*message.BroadcastTx).Tx, ShouldResemble, txRaw)
						So(broadcastRespMessage, ShouldHaveSameTypeAs, &message.BroadcastTxResponse{})
					})
				})
			})
		})
	})
}
