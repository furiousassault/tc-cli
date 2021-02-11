package token

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/furiousassault/tc-cli/pkg/configuration"
	"github.com/furiousassault/tc-cli/pkg/teamcity/subapi"
)

func CreateCommandTreeToken(config configuration.Configuration, serviceProvider serviceProvider) *cobra.Command {
	cmdToken := &cobra.Command{
		Use:   "token <subcommand>",
		Short: "token subcommand tree",
		RunE:  nil,
	}
	cmdTokenRotate := &cobra.Command{
		Use: "rotate <userID> <old_token_name> [<new_token_name>]",
		Short: "Revoke token with old token name " +
			"and create new token under the same name or under <new_token_name> if it is present. " +
			"New token will be written to token file path.",
		Args: validateTokenArgs,
		RunE: createHandlerTokenRotate(config, serviceProvider),
	}
	cmdToken.AddCommand(cmdTokenRotate)

	return cmdToken
}

type serviceProvider interface {
	TokenServiceCurrent() API
	TokenServiceWithTokenAuth(token string)
}

type API interface {
	TokenList(userID string) (tokenList subapi.Tokens, err error)
	TokenCreate(userID, tokenName string) (token subapi.Token, err error)
	TokenRemove(userID string, tokenName string) (err error)
}

type cobraCommandExecutorE func(_ *cobra.Command, args []string) error

func createHandlerTokenRotate(config configuration.Configuration, sp serviceProvider) cobraCommandExecutorE {
	return func(cmd *cobra.Command, args []string) error {
		userID := args[0]
		tokenNameOld := args[1]
		tokenNameNew := ""

		if len(args) > CommandTokenArgsNumberMin {
			tokenNameNew = args[2]
		}

		return tokenRotate(
			sp,
			config.API.Authorization.TokenFilePath,
			config.API.Authorization.Token,
			userID,
			tokenNameOld,
			tokenNameNew,
		)
	}
}

func tokenRotate(serviceProvider serviceProvider, tokenFilePath, tokenOld, userID, tokenNameOld, tokenNameNew string) error {
	tokens, err := serviceProvider.TokenServiceCurrent().TokenList(userID)
	if err != nil {
		return err
	}

	tokenNameTemporary := fmt.Sprint(time.Now().Unix())
	tokenTemporary, err := serviceProvider.TokenServiceCurrent().TokenCreate(userID, tokenNameTemporary)
	if err != nil {
		return err
	}

	// this should be performed in other way; refactor later
	serviceProvider.TokenServiceWithTokenAuth(tokenTemporary.Value)

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

	if err := serviceProvider.TokenServiceCurrent().TokenRemove(userID, tokenNameOld); err != nil {
		return err
	}

	if tokenNameNew == "" {
		tokenNameNew = tokenNameOld
	}

	token, err := serviceProvider.TokenServiceCurrent().TokenCreate(userID, tokenNameNew)
	if err != nil {
		return err
	}

	tokenWritePath := fmt.Sprintf("%s.%s.new", tokenFilePath, tokenNameNew)

	if err = ioutil.WriteFile(tokenWritePath, []byte(token.Value), 0600); err != nil {
		return err
	}

	if err := serviceProvider.TokenServiceCurrent().TokenRemove(userID, tokenNameTemporary); err != nil {
		return err
	}

	if err = os.Remove(tokenBackupFilePath); err != nil {
		return err
	}

	fmt.Printf(
		"Token with name '%s' has been rotated successfully. New token name is '%s'. Token file path: '%s'\n",
		tokenNameOld,
		tokenNameNew,
		tokenWritePath,
	)
	return nil
}
