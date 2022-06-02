# :vertical_traffic_light: zebra [![test](https://github.com/rchamarthy/zebra/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/rchamarthy/zebra/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/rchamarthy/zebra)](https://goreportcard.com/report/github.com/rchamarthy/zebra) [![codecov](https://codecov.io/gh/rchamarthy/zebra/branch/main/graph/badge.svg?token=94ZZ46W6VA)](https://codecov.io/gh/rchamarthy/zebra) [![Maintainability](https://api.codeclimate.com/v1/badges/e49a48cf30c0d644fd7b/maintainability)](https://codeclimate.com/github/rchamarthy/zebra/maintainability)

## Welcome to Zebra ##
Zebra is a tool to maintain resource inventory and reservations. 

### How Zebra works ###
Zebra is a neat and convenient tool for resource management. To start, any resource can be added to the system provided an ID string and other resource-specific details. Zebra must also be given the resource associations (i.e. how is the current resource connected to any other resources in the system). Once all resources and associations have been added, the inventory is complete. Zebra now models the entire system. From here, users can reserve the system resources. When a user reserves a resource, Zebra marks it as in-use. While the user holds the resource, Zebra continues to allocate free resources to subsequent users. When a user releases a resource, Zebra marks it as free.

### Zebra for Metrics ###
We aim to develop a dashboard to track resource usage by user.​ This provides insight on which resources are in high-demand, how each user is utilizing system resources, etc. A further enhancement would be to allow user groups. By doing so, Zebra can track usage across a user group and gain insight into how a group is using system resources.

### Getting started with Zebra ### 
As of now, we have not determined a way in which Zebra will read input data to model a system.

### Users ###
A user represents an temporary owner of a resource. Each user will be associated with a role. This role (such as developer, admin, client, etc.) determines the user's permissions. Once authenticated, a user will be allowed to reserve resources according to their role permissions. Once Zebra allocates a resource to the user, Zebra logs that the user is in current possession of the resource. Once the user is finished, Zebra will release the resource to be allocated to other users.
As of now, we have not determined how to create/delete/authenticate users to begin reserving resources.

### TO DO ###