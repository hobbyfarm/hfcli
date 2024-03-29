# hfcli: a simple cli to interact with hobbyfarm

hfcli is a simple cli to  make it easy to do certain operations on hobbyfarm.


```
/tmp/hfcli -h
Usage:
  hfcli [flags]
  hfcli [command]

Available Commands:
  apply       apply objects, valid options are scenario
  completion  generate the autocompletion script for the specified shell
  get         get objects, valid options are scenario  
  help        Help about any command
  info        perform info operations, valid options are accesscode and email

Flags:
  -h, --help                help for hfcli
  -k, --kubeconfig string   kubeconfig for authentication
  -n, --namespace string    namespace (default "gargantua")

Use "hfcli [command] --help" for more information about a command.
```

Currently, hfcli supports these tasks:

## apply

apply currently allows you to create a scenario by parsing a directory.

```
apply objects, valid options are scenario

Usage:
  hfcli apply [flags]
  hfcli apply [command]

Available Commands:
  scenario    create scenario NAME PATH_TO_SCENARIOS

Flags:
  -h, --help   help for apply

Global Flags:
  -k, --kubeconfig string   kubeconfig for authentication
  -n, --namespace string    namespace (default "gargantua")

Use "hfcli apply [command] --help" for more information about a command.
```

An example scenario is available in the `example` folder.

The folder needs to be structured as:


```
example
|   scenario.yml
|---content
|   |   step-1.md
|   |   step-2.md
|   |   step-n.md
```

Users can inject two metadata fields about the step into step as shown below.
```
 +++
 title = "heading for step 1"
 weight = 1
 +++
 
## Step 1
ls -lart
```

`title`: defines the name of the step, if not specified the name of the file is used as the step name

`weight`: defines the order of the step. Lower is set up earlier. If a step doesn't have a weight, then files are ordered alphabetically and added with a default weight.

## get

get downloads an existing scenario into a directory structure that can be used by apply.

```
get objects, valid options are scenario

Usage:
  hfcli get [flags]
  hfcli get [command]

Available Commands:
  scenario    get scenario NAME PATH_TO_SCENARIO

Flags:
  -h, --help   help for get

Global Flags:
  -k, --kubeconfig string   kubeconfig for authentication
  -n, --namespace string    namespace (default "gargantua")

Use "hfcli get [command] --help" for more information about a command.
```

## delete

deletes an existing scenario from hobbyfarm

```
deletes objects, valid options are scenario

Usage:
  hfcli deletes [flags]
  hfcli deletes [command]

Available Commands:
  scenario    deletes scenario NAME

Flags:
  -h, --help   help for get

Global Flags:
  -k, --kubeconfig string   kubeconfig for authentication
  -n, --namespace string    namespace (default "gargantua")

Use "hfcli get [command] --help" for more information about a command.
```

## info

info can be used to search for information about an accesscode or a user

```
perform info operations, valid options are accesscode and email

Usage:
  hfcli info [flags]
  hfcli info [command]

Available Commands:
  accesscode  hfcli info accesscode ACCESS_CODE
  email       get info about session and infra associated with email address

Flags:
  -h, --help   help for info

Global Flags:
  -k, --kubeconfig string   kubeconfig for authentication
  -n, --namespace string    namespace (default "gargantua")

Use "hfcli info [command] --help" for more information about a command.
```

`hfcli info accesscode CODE` can be used to search if an accesscode is in use.

If in use, the command will return all the current sessions in user with this accesscode

```
▶ /tmp/hfcli info accesscode hfcli
INFO[0001] scheduled event test has accesscode hfcli
SESSION        | VMID                       | STATUS       | PUBLICIP  |
ss-bvdn3nbytc  | dynamic-07583d7d-3b9ea64e  | provisioned  |           |
ss-bvdn3nbytc  | dynamic-07583d7d-7a405c86  | provisioned  |           |

```




`hfcli info email EMAILADDRESS` can be used to search for info about session, vm provisioning status for a specific user identified by the registration email address

```
▶ /tmp/hfcli info email admin
SESSION        | VMID                       | STATUS   | PUBLICIP     |
ss-bvdn3nbytc  | dynamic-07583d7d-3b9ea64e  | running  | 3.25.98.218  |
ss-bvdn3nbytc  | dynamic-07583d7d-7a405c86  | running  | 3.25.62.40   |
```