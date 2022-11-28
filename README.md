# :vertical_traffic_light: zebra [![test](https://github.com/project-safari/zebra/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/project-safari/zebra/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/project-safari/zebra)](https://goreportcard.com/report/github.com/project-safari/zebra) [![codecov](https://codecov.io/gh/project-safari/zebra/branch/main/graph/badge.svg?token=94ZZ46W6VA)](https://codecov.io/gh/project-safari/zebra) [![Maintainability](https://api.codeclimate.com/v1/badges/eb4a80c9a54fce474e28/maintainability)](https://codeclimate.com/github/project-safari/zebra/maintainability) [![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/project-safari/zebra/main/LICENSE)

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

### How to run Zebra ###
Zebra can be run by executing a simple script. This script can be found outside of this directory, in the parent directory, where the UI and backend are also situated.

* For a Linux/Ubuntu based system, simply navigate to that directory and use the following commands in terminal:

    1. chmod +x dev.sh

    2. run ./dev.sh

* For a Windows based system, follow the following steps*:

    1. Start -> Settings -> Update&Security. Under the ‘Use Developer Features’, select ‘Developer mode’.

    2. Selecting the developer mode will pop an alert. Click yes, and let the computer restart.

    3. Go to Control Panel -> Programs and Features -> Turn Windows Features On and Off. In the window that appears, check the ‘Windows Subsystem for Linux’ option, and click OK.

    4. This will trigger an alert asking for the system to be restarted to complete the installation of the required components. After the restart is complete, go to the command prompt, and type ‘bash’. Follow the instructions that appear to install bash from Windows store. After it is installed, it will be required to create a UNIX username. After completing the installation, exit the prompt.

    5. To access the shell, simply type ‘bash’ in the Windows command prompt, and everything is good to go.

    *Steps work with Windows 10 and above. Additionally, the 64-bit bersion of the OS is needed. Note that for previous Windows versions, emulators like Cygwin or Git need to be installed on the host machine.

### TO DO ###