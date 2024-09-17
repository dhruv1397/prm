package purge

import (
	"bufio"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/dhruv1397/prm/cli"
	"github.com/dhruv1397/prm/store"
	"os"
	"strings"
)

type command struct {
	force bool
}

func (c *command) run(*kingpin.ParseContext) error {
	if !c.force {
		confirmation := promptForConfirmation()
		if !confirmation {
			fmt.Println("Purge aborted.")
			return nil
		}
		fmt.Println("Proceeding with the purge...")
	}
	str := store.NewSCMProviderImpl()
	err := str.Purge()
	if err != nil {
		return err
	}
	fmt.Println("Purge completed.")
	return nil
}

func Register(app *kingpin.Application) {
	c := &command{}

	cmd := app.Command(cli.CommandPurge, cli.CommandPurgeHelpText).Action(c.run)
	cmd.Flag(cli.FlagForce, cli.FlagForceHelpText).Short(cli.FlagOutputForce).BoolVar(&c.force)
}

func promptForConfirmation() bool {
	fmt.Print("This will delete all the providers, do you want to continue? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	userInput = strings.TrimSpace(userInput)
	return strings.ToLower(userInput) == "y" || strings.ToLower(userInput) == "yes"
}
