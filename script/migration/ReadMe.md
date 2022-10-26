            ReadMe about Migration of Data from the Racktables Database

This ReadMe file contains information about the structure of the databases and about how migration works.

MigrationFromRacktables.go is a script to provide data migration from the racktables database to the Zebra tool.

The database that is used for this script is a MariaDB database. The user should be advised that the MariaDB mysql user utilized here has the following credentials: 
                    
                        username: username
                        password: 1234

    The racktables database has several tables which contain information about resources that are relevant to the Zebra tool.

                        +---------------------------+
                        | Tables_in_racktables      |
                        +---------------------------+
                        | Atom                      |
                        | Attribute                 |
                        | AttributeMap              |
                        | AttributeValue            |
                        | CachedPAV                 |
                        | CachedPNV                 |
                        | CachedPVM                 |
                        | CactiGraph                |
                        | CactiServer               |
                        | Chapter                   |
                        | Config                    |
                        | Dictionary                |
                        | EntityLink                |
                        | File                      |
                        | FileLink                  |
                        | IPv4Address               |
                        | IPv4Allocation            |
                        | IPv4LB                    |
                        | IPv4Log                   |
                        | IPv4NAT                   |
                        | IPv4Network               |
                        | IPv4RS                    |
                        | IPv4RSPool                |
                        | IPv4VS                    |
                        | IPv6Address               |
                        | IPv6Allocation            |
                        | IPv6Log                   |
                        | IPv6Network               |
                        | LDAPCache                 |
                        | Link                      |
                        | Molecule                  |
                        | MountOperation            |
                        | MuninGraph                |
                        | MuninServer               |
                        | Object                    |
                        | ObjectHistory             |
                        | ObjectLog                 |
                        | ObjectParentCompat        |
                        | PatchCableConnector       |
                        | PatchCableConnectorCompat |
                        | PatchCableHeap            |
                        | PatchCableHeapLog         |
                        | PatchCableOIFCompat       |
                        | PatchCableType            |
                        | Port                      |
                        | PortAllowedVLAN           |
                        | PortCompat                |
                        | PortInnerInterface        |
                        | PortInterfaceCompat       |
                        | PortLog                   |
                        | PortNativeVLAN            |
                        | PortOuterInterface        |
                        | PortVLANMode              |
                        | RackSpace                 |
                        | RackThumbnail             |
                        | Script                    |
                        | TagStorage                |
                        | TagTree                   |
                        | UserAccount               |
                        | UserConfig                |
                        | VLANDescription           |
                        | VLANDomain                |
                        | VLANIPv4                  |
                        | VLANIPv6                  |
                        | VLANSTRule                |
                        | VLANSwitch                |
                        | VLANSwitchTemplate        |
                        | VLANValidID               |
                        | VS                        |
                        | VSEnabledIPs              |
                        | VSEnabledPorts            |
                        | VSIPs                     |
                        | VSPorts                   |
                        | location                  |
                        | rack                      |
                        | rackobject                |
                        | row                       |
                        +---------------------------+

        This migration script is particularly concerned with the following tables:

        1. Rackspace

            This table contains information about each rack with its respective contents. 
                
                This table's columns are:

                    a. rack_id, used to identify each respective rack

                    b. object_id, used to identify each resource in a given rack and subsequently make queries about it.  
        
        2. IPv4Allocation

            Some resources in the Zebra tool require an IP address. This table contains IP addresses and object IDs for resources. 
            
                This table's columns, which are used in the script, are:

                    1. ip, used to provide certain resources with the necessary IP addresses, as well as to further make queries based on IP address for further information, such as owner information. 

                    2. object_id, used to identify each individual's resource IP information.

        3. Port

            Some resources in the Zebra tool contain ports. This table contains port information for each individual resource. The Zebra tool is interested in finding the number of ports in each resource. Some resources only have one port, others have multiple ports. Thus, for each resource in the database, this number is calculated and assigned. 
                
                The column of interest in this table is:

                    1. id, which represents the ID information of each port rather than the resource's ID. Each port has its unique ID number. 

        4. Rackobject

            Each rack in the Racktables database contains several resources, hereby named rackobjects. Each of these resources has additional information to that mentioned above. This information will be assigned to resources in the Zebra tool, as needed.

                The columns of interest from this table are:

                    1. id, which reresents each individual's unique ID.

                    2. name, which represents each individual's resource name. This is different than resource's type.

                    3.label, which represents eiter a key-value pair, a "tag," or a system.group label that identifes the group that a resource belongs to.

                    4. objtype_id, which is a numerical value that represents a resource type. This numerical value is the same for a given group of resources, such as servers. This numerical value is used to determine the type.

                    5. asset_no, which represents the number of assets of each resource.

                    6. has_problems, which is a yes or no value about a resource's state.

                    7. comment, which represents any additional notes about a resource.

        5. Rack

            The Rack table contains additional information about each rack in the database and about the respective rows.

                The columns used from this table are:

                    1. row_id, which represents the ID number of each row, given a rack ID.

                    2. row_name, which represents the name of each row, given a rack ID.

                    3. location_name, which is ina letter and number string format that represents the location of a rack and row. 

        6. IPv4Log

            Among other pieces of data, this table contains information about users, given a particular IP address.

                The column of interest from this table is:

                1. user, which is the name / role of a user that owns / manages a particular resource.             
                
            ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~