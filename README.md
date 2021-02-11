# tc-cli

Basic console client for TeamCity.

For now, it supports a limited number of scenarios. API objects are partially deserialized, with support of the
currently required fields only. Should not be used as "library", package, etc. Backward compatibility is not guaranteed.

Tested only with current latest version of Teamcity (2020.2.1 (build 85633)), Supports token-based, HTTP and guest
authentication. Token-based authentication is encouraged.

### Configuration

The application reads configuration from YAML file which must be specified via `-c/--configpath` with default value
`~/.tc-client/configuration.yaml`.

Default may be overridden by setting `TC_CLI_CONFIG_PATH` environment variable. No default configuration is provided at
this moment.

### Build

Application can be built using `go build cmd/main.go`. Another quick starter is

```
make build-image
docker run -e TC_CLI_TOKEN=<your_token> --network=host furiousassault/tc-cli:latest -c configuration.example.yaml list projects
```

provided, that configuration file contains valid URL of Teamcity server.

### Supported scenarios

#### List projects

```
$tc-cli -c configuration.example.yaml list projects
ID              NAME                    DESCRIPTION                 
_Root           "<Root project>"        Contains all other projects     
TestProject0    "Test Project 0"        Yet another project
```

#### List build configurations

```
$tc-cli -c configuration.example.yaml list buildtypes TestProject0
ID                      NAME           
TestProject0_BuildConf1 "Build conf 1"  
TestProject0_BuildConf0 "BuildConf0"
```

#### List builds of build configuration

```
$tc-cli -c configuration.example.yaml list builds TestProject0_BuildConf0
ID      NUMBER  STATE           STATUS  
202     14      finished        SUCCESS 
201     13      finished        SUCCESS 
109     12      finished        SUCCESS
```

The command accepts optional argument `-n/--number` which limits count of returned builds to N most recent.

#### Describe build

The command shows input and output parameters of build specified by build configuration ID and build number.

```
$tc-cli -c configuration.example.yaml describe build TestProject0_BuildConf0 14
ID      NUMBER  STATE           STATUS  QUEUED                          STARTED                         FINISHED                      
202     14      finished        SUCCESS 2021-02-07 08:44:48 +0000 UTC   2021-02-07 08:44:50 +0000 UTC   2021-02-07 08:44:55 +0000 UTC   

Properties
KEY     VALUE           INHERITED 
TESTENV testenvvalue    true

Resulting properties
KEY                                                             VALUE                                                                                                   
build.counter                                                   14                                                                                                      
[...]                                                           
DotNetCredentialProvider1.0.0_Path                              /opt/buildagent/plugins/nuget-agent/bin/credential-plugin/netcoreapp1.0/CredentialProvider.TeamCity.dll 
```

The output of resulting properties, which might contain long list of variables, can be suppressed by
setting `-s/--short` flag.

#### Build log

Log of particular finished build can be obtained by `log` command, with build ID specified.

```
$tc-cli -c configuration.example.yaml log 202
Build 'Test Project 0 / BuildConf0' #14
...
Started 2021-02-07 08:44:50 on agent 'ip_172.17.0.1'
Finished 2021-02-07 08:44:55 with status NORMAL 'Success'
[08:44:48] : bt1 (6s)
[08:44:48]i: TeamCity server version is 2020.2.1 (build 85633)
[08:44:48] : The build is removed from the queue to be prepared for the start
[08:44:48] : Collecting changes in 1 VCS root (1s)
...
```

Some inconsistency about using ids and numbers of builds in different commands takes place for now, should be fixed in
future. Log is printed to stdout for now, like in `kubectl logs` default behavior, output file support should be added
further.

#### Run build configuration

```
$tc-cli -c configuration.example.yaml run TestProject0_BuildConf0
BUILD ID        QUEUED BY               STATE  
203             furiousassault        queued  

Properties
TYPE    KEY     VALUE           INHERITED 
env     TESTENV testenvvalue    true 
```

Custom build parameters or build comments are not supported yet.

#### Rotate token

Usage:

```
token rotate <userID> <old_token_name> [<new_token_name>]`
```

Example:

```
$tc-cli -c configuration.example.yaml token rotate furiousassault token_0
Token with name 'token_0' has been rotated successfully. New token name is 'token_0'. Token file path: '/home/furiousassault/workspace/tc-cli/access.token.token_0.new'

$tc-cli token -c configuration.example.yaml rotate furiousassault token_1 token_11
Token with name 'token_1' has been rotated successfully. New token name is 'token_11'. Token file path: '/home/furiousassault/workspace/tc-cli/access.token.token_11.new'
```

Unfortunately, there's no obvious way to get the token name by its value nor to get userID by token, so user has to
specify his userID, and the name of token that is going to be revoked.

User can optionally specify new token name as the last parameter. New token is stored in
the `token_file_path.<token_new>.new` path specified. `<token_new>` defaults to old token name if `new_token_name` is
omitted.

This is done to prevent the unwanted override, if the token user wants to rotate is not his current work token, with
which application is being run.

When the better way to perform such operation will be discovered (in API or by providing more complex token local
storage), the behavior will change.

#### Download artifact

Usage:

```
artifact download <buildID> <path>
```

Example:

```
$tc-cli -c configuration.example.yaml artifact download 202 tmpartifacts/step1.art
Artifact has been downloaded to '/tmp/tc-client/artifacts/202/tmpartifacts/step1.art'

$tc-cli -c configuration.example.yaml artifact download 202 tmpartifacts/step1.art -o /tmp/withoutprefix.step1.art
Artifact has been downloaded to '/tmp/withoutprefix.step1.art'

```

Output path may be specified by `-o/--output` flag. The value is interpreted as directory with adding path
suffix `<buildID>/<artifact_path>`, if it ends with slash, dot or double dot. An attempt to create nested directories
will be done. Otherwise, it's interpreted as direct path to write. No path suffixes will be added in such case.

Default artifacts directory path is `/tmp/tc-client/artifacts/`. Default can be overridden in configuration by
setting `artifacts_directory`
parameter, or by environment variable `TC_CLI_ARTIFACTS_DIRECTORY_DEFAULT`.

If the target output file exists, command requires flag `-f/--force` to override it.
