package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gantsign/alt-galaxy/internal/roleinstaller"
	"gopkg.in/urfave/cli.v1"
)

var (
	version  = "unknown"
	revision = "unknown"
	built    = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Alternate implementation of ansible-galaxy tool for downloading Ansible roles."
	app.HideHelp = true

	app.Version = version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("Version:  ", version)
		fmt.Println("Revision: ", revision)
		fmt.Println("Built:    ", built)
	}

	app.Commands = []cli.Command{
		{
			Name:      "install",
			Usage:     "Install Ansible roles from the specified role file.",
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "role-file",

					Usage: "A file containing a list of roles to be imported",
				},
				cli.StringFlag{
					Name:  "roles-path",
					Usage: "The path to the directory containing your roles. The default is the roles_path configured in your ansible.cfg file (/etc/ansible/roles if not configured",
				},
			},
			Action: func(c *cli.Context) error {
				roleFile := c.String("role-file")
				if roleFile == "" {
					return errors.New("You must specify a role file.")
				}

				rolesPath := c.String("roles-path")
				if rolesPath == "" {
					return errors.New("You must specify the roles path.")
				}

				cmd := roleinstaller.RoleInstallerCmd{
					RoleFile:  roleFile,
					RolesPath: rolesPath,
					ApiServer: "https://galaxy.ansible.com",
				}
				return cmd.Execute()
			},
		},
	}

	appError := app.Run(os.Args)
	if appError != nil {
		fmt.Fprintln(os.Stderr, appError)
		os.Exit(1)
	}
}
