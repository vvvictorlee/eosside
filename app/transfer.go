package app

import (

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/blockchain-develop/eosside/x/ibc"
)

const (
	flagEosAccount   = "eos-account"
	flagAmount = "amount"
)

// IBC relay command
func SideTransferCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			// get the from address
			from, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}

			// build the message
			msg, err := buildSideTransferMsg(from)
			if err != nil {
				return err
			}

			// get password
			err = ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, []sdk.Msg{msg}, cdc)
			if err != nil {
				return err
			}
			return nil
		},
	}
	
	cmd.Flags().String(flagEosAccount, "", "eos account to transfer")
	cmd.MarkFlagRequired(flagEosAccount)
	viper.BindPFlag(flagEosAccount, cmd.Flags().Lookup(flagEosAccount))

	cmd.Flags().String(flagAmount, "", "transfer amount")
	cmd.MarkFlagRequired(flagAmount)
	viper.BindPFlag(flagAmount, cmd.Flags().Lookup(flagAmount))
	
	return cmd
}

func buildSideTransferMsg(from sdk.AccAddress) (sdk.Msg, error) {
	//
	amount := viper.GetString(flagAmount)
	coins, err := sdk.ParseCoins(amount)
	if err != nil {
		return nil, err
	}
	
	//
	msg := ibc.SideTransferMsg {
		SrcAddr: from,
		Coins: coins,
		DestAddr: viper.GetString(flagEosAccount),
	}

	return msg, nil
}
