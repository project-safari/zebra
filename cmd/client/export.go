package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
	"github.com/spf13/cobra"
)

// put store into json file

var ErrExport = errors.New("error with export command")

func NewExport() *cobra.Command {
	importCmd := &cobra.Command{
		Use:          "export {file name}.json",
		Short:        "export resources into given json file",
		RunE:         exportResources,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
	}

	return importCmd
}

func exportResources(cmd *cobra.Command, args []string) error {
	cfg, req, err := makeExportReq(cmd, args)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	res := zebra.NewResourceMap(store.DefaultFactory())

	resCode, err := client.Get("api/v1/resources", req, res)
	if resCode != http.StatusOK {
		return ErrExport
	}

	if err != nil {
		return err
	}

	out, err := os.Create(args[0])
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}

	if _, err := out.Write(bytes); err != nil {
		return err
	}

	if err := out.Close(); err != nil {
		return err
	}

	fmt.Printf("Exported resources to %s.\n", args[0])

	return err
}

func makeExportReq(cmd *cobra.Command, args []string) (*Config, interface{}, error) {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, err := Load(cfgFile)
	if err != nil {
		return nil, nil, err
	}

	// query for all resources
	qr := struct {
		IDs        []string      `json:"ids,omitempty"`
		Types      []string      `json:"types,omitempty"`
		Labels     []zebra.Query `json:"labels,omitempty"`
		Properties []zebra.Query `json:"properties,omitempty"`
	}{
		Types: []string{},
	}

	return cfg, qr, nil
}
