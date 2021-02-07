package token

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/configuration"
	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

func CreateCommandTreeToken(config configuration.Configuration, tokenAPI API) *cobra.Command {
	cmdToken := &cobra.Command{
		Use:   "token <subcommand>",
		Short: "token subcommand tree",
		RunE:  nil,
	}
	cmdTokenRotate := &cobra.Command{
		Use:   "rotate <userID> <old_token_name> <new_token_name>",
		Short: "show main attributes of entity specified by id",
		Args:  validateTokenArgs,
		RunE:  createHandlerTokenRotate(config, tokenAPI),
	}
	cmdToken.AddCommand(cmdTokenRotate)

	return cmdToken
}

type API interface {
	TokenList(userID string) (tokenList subapi.Tokens, err error)
	TokenCreate(userID, tokenName string) (token subapi.Token, err error)
	TokenRemove(userID string, tokenName string) (err error)
}

type cobraCommandExecutorE func(_ *cobra.Command, args []string) error

func createHandlerTokenRotate(config configuration.Configuration, tokenAPI API) cobraCommandExecutorE {
	return func(cmd *cobra.Command, args []string) error {
		userID := args[0]
		tokenNameOld := args[1]
		tokenNameNew := args[2]

		return tokenRotate(
			tokenAPI,
			config.API.Authorization.TokenFilePath,
			config.API.Authorization.Token,
			userID,
			tokenNameOld,
			tokenNameNew,
		)
	}
}

func tokenRotate(tokenAPI API, tokenFilePath, tokenOld, userID, tokenNameOld, tokenNameNew string) error {
	tokens, err := tokenAPI.TokenList(userID)
	if err != nil {
		return err
	}

	for _, token := range tokens.Items {
		if tokenNameNew == token.Name {
			return errors.Wrapf(errTokenNamesMatch, "token with name \"%s\" already exists", tokenNameNew)
		}
	}

	// create new token using user credentials
	tokenBackupFilePath := fmt.Sprintf("%s.old", tokenFilePath)
	if err := ioutil.WriteFile(tokenBackupFilePath, []byte(tokenOld), 0600); err != nil {
		return err
	}

	token, err := tokenAPI.TokenCreate(userID, tokenNameNew)
	if err != nil {
		return err
	}

	if err := tokenAPI.TokenRemove(userID, tokenNameOld); err != nil {
		return err
	}

	if err = ioutil.WriteFile(tokenFilePath, []byte(token.Value), 060); err != nil {
		return err
	}

	if err = os.Remove(tokenBackupFilePath); err != nil {
		return err
	}

	fmt.Printf("Token with name '%s' has been rotated successfully.\n", tokenNameOld)
	return nil
}
