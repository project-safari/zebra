package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/project-safari/zebra"
	"github.com/spf13/cobra"
)

var ErrShow = errors.New("error with show command")

func NewShow() *cobra.Command {
	showCmd := &cobra.Command{
		Use:          "show",
		Short:        "show resources (filter by type)",
		RunE:         showResources,
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(0),
	}

	showCmd.Flags().StringP("type", "t", "all", "resource type")

	return showCmd
}

func showResources(cmd *cobra.Command, args []string) error {
	cfg, req, err := makeShowReq(cmd, args)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	res := zebra.NewResourceMap(nil)

	resCode, err := client.Get("api/v1/resources", req, res)
	if resCode != http.StatusOK {
		return ErrShow
	}

	fmt.Println(res.MarshalJSON())

	return err
}

func makeShowReq(cmd *cobra.Command, args []string) (*Config, interface{}, error) {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, err := Load(cfgFile)
	if err != nil {
		return nil, nil, err
	}

	types, err := cmd.Flags().GetStringSlice("type")
	if err != nil {
		return nil, nil, err
	}

	// If all, make an empty request to show all resources.
	if types[0] == "all" {
		types = []string{}
	}

	qr := struct {
		IDs        []string      `json:"ids,omitempty"`
		Types      []string      `json:"types,omitempty"`
		Labels     []zebra.Query `json:"labels,omitempty"`
		Properties []zebra.Query `json:"properties,omitempty"`
	}{
		Types: types,
	}

	return cfg, qr, nil
}
