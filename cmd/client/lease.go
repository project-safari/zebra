package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/lease"
	"github.com/spf13/cobra"
)

var ErrCreateLease = errors.New("error creating resource")

const DefaultResourceCount = 3

func NewLease() *cobra.Command {
	leaseCmd := &cobra.Command{
		Use:          "lease",
		Short:        "request a lease",
		RunE:         leaseRequest,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
	}

	leaseCmd.Flags().StringP("group", "g", "global", "resource group")
	leaseCmd.Flags().IntP("count", "k", DefaultResourceCount, "number of resources")

	return leaseCmd
}

func leaseRequest(cmd *cobra.Command, args []string) error {
	cfg, req, resReq, err := makeLeaseReq(cmd, args)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	// Create a new lease
	resCode, err := client.Post("api/v1/resources", req, nil)
	if resCode != http.StatusOK {
		return ErrCreateLease
	}

	fmt.Println("Request - Type:", resReq.Type, "Group:", resReq.Group, "Count:", resReq.Count)
	fmt.Println("Lease request successfully created")

	return err
}

func makeLeaseReq(cmd *cobra.Command, args []string) (*Config, *zebra.ResourceMap, *lease.ResourceReq, error) {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, err := Load(cfgFile)
	if err != nil {
		return nil, nil, nil, err
	}

	resCount, err := cmd.Flags().GetInt("count")
	if err != nil {
		return nil, nil, nil, err
	}

	req := &lease.ResourceReq{
		Type:  args[0],
		Group: cmd.Flag("group").Value.String(),
		Count: resCount,
	}

	lease := lease.NewLease(
		cfg.Email,
		time.Duration(cfg.Defaults.Duration)*time.Hour,
		[]*lease.ResourceReq{req})

	resMap := zebra.NewResourceMap(model.Factory())
	if err := resMap.Add(lease); err != nil {
		return nil, nil, nil, err
	}

	return cfg, resMap, req, nil
}
