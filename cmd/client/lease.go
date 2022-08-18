package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/lease"
	"github.com/project-safari/zebra/store"
	"github.com/spf13/cobra"
)

// ErrCreateLease returns an error message if a request.
//
// Used for creating a lease on a certain resource fails.
var ErrCreateLease = errors.New("error creating resource")

// The default number of resources, currently set to 3.
const DefaultResourceCount = 3

// Command for new lease(s).
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

// Creating a request for a new lease.
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

// Function to complete the lease request.
//
// Uses the cobra command and string arguments.
//
// Loads the config file, gets the count, resource type, and resource group.
//
// Helps get the new lease according to the request.
//
// Returns pointers to the config, a resource map, a lease resource request, and (a) potential error(s).
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

	resMap := zebra.NewResourceMap(store.DefaultFactory())
	resMap.Add(lease, lease.GetType())

	return cfg, resMap, req, nil
}
