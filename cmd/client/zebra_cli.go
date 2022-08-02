package main

import (
	"github.com/spf13/cobra"
)

/*
Build a zebra inventory command line client:
zebra show servers
zebra show users
zebra show registrations
zebra show networks
zebra show switches
zebra show racks
zebra show labs
zebra show data-centers
zebra command will use the private RSA key that the user created
before registration and authenticate so there is no need to login for the CLI
all show commands will support label filters

distributed among the following files:
zebra_cli (command inventory)
run_cli (runE functions)
show_test (test file)
*/

// create inventory of commands for the zebra cli.
func NewZebra() *cobra.Command {
	// default zebra command to show resources.
	zebraCmd := &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:   "show",
		Short: "show resources",
	}

	usrCmd := &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "users",
		Short:        "show zebra users",
		RunE:         ShowUsr,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	}

	usrCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "registrations",
		Short:        "show zebra registrations",
		RunE:         ShowReg,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	zebraCmd.AddCommand(usrCmd)

	return zebraCmd
}

// function to add show commands for network resources to zebraCmd.
func NewNetCmd(zebraCmd *cobra.Command) *cobra.Command {
	zebraCmd.AddCommand(
		&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
			Use:          "vlans",
			Short:        "show zebra switches",
			RunE:         ShowVlan,
			Args:         cobra.MaximumNArgs(1),
			SilenceUsage: true,
		})

	zebraCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "sws",
		Short:        "show zebra switches",
		RunE:         ShowSw,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	zebraCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "ip",
		Short:        "show zebra IP-Address-Pools",
		RunE:         ShowIP,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	return zebraCmd
}

// function to add show commands for dc resources to zebraCmd.
func NewDCCmd(zebraCmd *cobra.Command) *cobra.Command {
	zebraCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "dc",
		Short:        "show datacenters",
		RunE:         ShowDC,
		SilenceUsage: true,
	})

	zebraCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "labs",
		Short:        "show labs",
		RunE:         ShowLab,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	zebraCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "racks",
		Short:        "show racks",
		RunE:         ShowRack,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	return zebraCmd
}

// function to add show commands for server resources to zebraCmd.
func NewSrvCmd(zebraCmd *cobra.Command) *cobra.Command {
	zebraCmd.AddCommand(&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "servers",
		Short:        "show servers",
		RunE:         ShowServ,
		SilenceUsage: true,
	})

	zebraCmd.AddCommand((&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "esx",
		Short:        "show esx-servers",
		RunE:         ShowESX,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	}))

	zebraCmd.AddCommand((&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "vcenters",
		Short:        "show vcenters",
		RunE:         ShowVC,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	}))

	zebraCmd.AddCommand((&cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "vms",
		Short:        "show vms",
		RunE:         ShowVM,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	}))

	return zebraCmd
}
