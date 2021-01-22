package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/configuration"
	"github.com/furiousassault/tc-cli/pkg/teamcity/api"
	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

var (
	cmdToken = &cobra.Command{
		Use:   "token <subcommand>",
		Short: "token subcommand tree",
		RunE:  nil,
	}

	cmdTokenRotate = &cobra.Command{
		Use:   "rotate <userID> <old_token_name> <new_token_name>",
		Short: "show main attributes of entity specified by id",
		Args:  validateTokenArgs,
		RunE:  tokenRotator(nil, tokenRotate),
	}
)

type TokenAPI interface {
	TokenLister
	TokenCreator
	TokenRemover
}

type TokenLister interface {
	TokenList(userID string) (tokenList subapi.Tokens, err error)
}

type TokenCreator interface {
	TokenCreate(userID, tokenName string) (token subapi.Token, err error)
}

type TokenRemover interface {
	TokenRemove(userID string, tokenName string) (err error)
}

type tokenRotateFunc func(tokenAPI TokenAPI, _ *cobra.Command, args []string) error
type cobraCommandExecutorE func(_ *cobra.Command, args []string) error

func tokenRotator(tokenAPI TokenAPI, trf tokenRotateFunc) cobraCommandExecutorE {
	return func(cmd *cobra.Command, args []string) error {
		// FIXME GLOBAL VARIABLE WRONG INIT
		return trf(api.API().Token, cmd, args)
	}
}

func tokenRotate(tokenAPI TokenAPI, cmd *cobra.Command, args []string) error {
	tokenFilePath := configuration.GetConfig().API.Authorization.TokenFilePath
	tokenOld := configuration.GetConfig().API.Authorization.Token

	tokens, err := tokenAPI.TokenList(args[0])
	if err != nil {
		return err
	}

	for _, token := range tokens.Items {
		if args[2] == token.Name {
			return fmt.Errorf("token with name \"%s\" already exists", args[2])
		}
	}

	// create new token using user credentials
	tokenBackupFilePath := fmt.Sprintf("%s.old", tokenFilePath)
	if err := ioutil.WriteFile(tokenBackupFilePath, []byte(tokenOld), 0666); err != nil {
		return err
	}

	token, err := tokenAPI.TokenCreate(args[0], args[2])
	if err != nil {
		return err
	}

	if err := tokenAPI.TokenRemove(args[0], args[1]); err != nil {
		return err
	}

	if err = ioutil.WriteFile(tokenFilePath, []byte(token.Value), 0666); err != nil {
		return err
	}

	if err = os.Remove(tokenBackupFilePath); err != nil {
		return err
	}

	fmt.Printf("Token with name '%s' has been rotated successfully.\n", args[1])
	return nil
}

func validateTokenArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(3)(cmd, args); err != nil {
		return err
	}

	if args[1] == args[2] {
		return errors.New("old and new token names must not be equal")
	}

	return nil
}
