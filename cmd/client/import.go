package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
	"github.com/spf13/cobra"
)

// should import from a file -- do a bulk create
// one arg, file name
// file will have a resource map

var ErrImport = errors.New("error with import command")

func NewImport() *cobra.Command {
	importCmd := &cobra.Command{
		Use:          "import {file name}.json",
		Short:        "import resources given json file",
		RunE:         importResources,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
	}

	return importCmd
}

func importResources(cmd *cobra.Command, args []string) error {
	cfg, req, err := makeImportReq(cmd, args)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	resCode, err := client.Post("api/v1/resources", req, nil)
	if resCode != http.StatusOK {
		return ErrImport
	}

	if err != nil {
		return err
	}

	fmt.Printf("Imported resources from %s.\n", args[0])

	return err
}

func makeImportReq(cmd *cobra.Command, args []string) (*Config, interface{}, error) {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, err := Load(cfgFile)
	if err != nil {
		return nil, nil, err
	}

	req := zebra.NewResourceMap(store.DefaultFactory())

	jsonFile, err := os.Open(args[0])
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err := json.Unmarshal(byteValue, &req); err != nil {
		return nil, nil, err
	}

	return cfg, req, nil
}
