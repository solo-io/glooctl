       ________                 __  __
      / ____/ /___  ____  _____/ /_/ /
     / / __/ / __ \/ __ \/ ___/ __/ / 
    / /_/ / / /_/ / /_/ / /__/ /_/ /  
    \____/_/\____/\____/\___/\__/_/   

CLI for Gloo

## Introduction
`glooctl` is command line tool to manage Gloo resources like Upstream and Virtual Host.
It provides you a friendlier way to manage the Gloo resources.

## Getting Started
Download the latest release of `glooctl` from https://github.com/solo-io/glooctl/releases/latest/

If you prefer to compile your own binary or work on the development of `glooctl` please use
the following command:


```
go get github.com/solo-io/glooctl
```

All `glooctl` commands take `kubeconfig` and `namespace` parameters. If these are not provided,
they default to `~/.kube/config` and `gloo-system` respectively.

### Map a Route to a Function

```
glooctl route map --path-prefix /calculator --function aws_lambda:calc
```
Uses default virtual host (named 'default'?)
should use -V to specify virtual host or can we use --host hhh to find the virtual host or should?
should we create a virtual host if not present?

will create route even if function doesn't exist

## Managing Routes on a Virtual Host
The `route` command allows you to manage the routes on a specific
virtual host. All `route` subcommands take the `vhost` or `V` flag
to specify the virtual host you want to manage.

### Getting Routes
The `get` command returns a list of routes on the virtual host.

```
glooctl route get -V vhost1

request exact path: /bar
request path prefix: /foo
event matcher: /apple
```

By default, `get` returns a summary list. You can pass the `output`
flag to see response in YAML or JSON to get details of the routes.

```
glooctl route get -V vhost1 -o yaml

extensions:
  auth:
    credentials:
      Password: bob
      Username: alice
    token: my-12345
request_matcher:
  path_exact: /bar
  verbs:
  - GET
  - POST
single_destination:
  upstream:
    name: my-upstream

extensions:
  auth:
    credentials:
      Password: bob
      Username: alice
    token: my-12345
request_matcher:
  headers:
    x-foo-bar: ""
  path_prefix: /foo
  verbs:
  - GET
  - POST
single_destination:
  function:
    function_name: foo
    upstream_name: aws

event_matcher:
  event_type: /apple
extensions:
  auth:
    credentials:
      Password: bob
      Username: alice
    token: my-12345
single_destination:
  function:
    function_name: foo
    upstream_name: aws
```

### Deleting a Route

```
glooctl route delete -V vhost1 --path-prefix /foo

request exact path: /bar
event matcher: /apple
```
### Appending a new Route

```
request_matcher:
  path_prefix: /foo/bar
  verbs:
  - GET
  - POST
single_destination:
    upstream:
      name: upstream2
```

```
glooctl route append -V vhost1 -f route.yaml 
request exact path: /bar
event matcher: /apple
request path prefix: /foo/bar
```

### Sorting Routes

```
glooctl route sort -V vhost1  
event matcher: /apple
request exact path: /bar
request path prefix: /foo/bar
```


## Managing Upstreams
`glooctl` provides a manual method of managing Upstreams. Gloo provides auto discovery 
service that can create or delete upstreams automatically. It also provides function
discovery service to manage the functions in an Upstream.

### Creating Upstream
The CLI allows you to create an upstream from a YAML file. 

Let's look at an upstream definition in `upstream.yaml`

```
name: testupstream
type: aws
spec:
  region: "us-east-2"
  credential: "aws-secret"
```

If you want to see the newly created upstream, you can pass `output` flag.

```
glooctl upstream create -f upstream.yaml --output yaml

Upstream created
metadata:
  namespace: gloo-system
  resource_version: "224352"
name: testupstream
spec:
  credential: aws-secret
  region: us-east-2
type: aws
```

### Getting Upstream
By default, `get` command returns a list of upstream names. 

```
glooctl upstream get

testupstream
```

You can pass it the `output` flag to return it as JSON or YAML.

```
glooctl upstream get -o json

{"name":"testupstream","type":"aws","spec":{"credential":"aws-secret","region":"us-east-2"},"metadata":{"resource_version":"224352","namespace":"gloo-system"}}
```

If you want to get details of a specific Upstream, you can use
`get` command with the name of the upstream. It returns
the result as YAML, but you can use `output` flag to get JSON.

```
glooctl upstream get testupstream

metadata:
  namespace: gloo-system
  resource_version: "224352"
name: testupstream
spec:
  credential: aws-secret
  region: us-east-2
type: aws
```

### Updating Upstream
Similar to `create` command, `update` command takes the definition of
upstream from a YAML file and replaces the existing upstream with the
one from the file.

```
glooctl upstream update -f upstream2.yaml -o yaml

Upstream updated
metadata:
  namespace: gloo-system
  resource_version: "224867"
name: testupstream
spec:
  credential: aws-secret
  region: us-east-1
type: aws
```

### Deleting Upstream
You can delete an existing upstream by giving the name of the
upstream to be deleted to `delete` command.

```
glooctl upstream delete testupstream

Upstream testupstream deleted
```

## Managing Virtual Hosts
`glooctl` provides a manual method of managing Virtual Hosts. Gloo provides auto discovery 
service that can create or delete virtual hosts automatically. 

### Creating Virtual Host
The CLI allows you to create a virtual from a YAML file. 

Let's look at a virtual host definition in `vhost.yaml`

```
name: vhost1
routes:
- request_matcher:
    path_exact: /bar
    verbs:
    - GET
    - POST
  single_destination:
    upstream:
      name: my-upstream
```

If you want to see the newly created virtual, you can pass `output` flag.

```
glooctl vhost create -f vhost.yaml --output yaml

Virtual host created  vhost1
metadata:
  namespace: gloo-system
  resource_version: "226902"
name: vhost1
routes:
- request_matcher:
    path_exact: /bar
    verbs:
    - GET
    - POST
  single_destination:
    upstream:
      name: my-upstream
```

### Getting Virtual Host
By default, `get` command returns a list of virtual host names. 

```
glooctl vhost get

vhost1
```

You can pass it the `output` flag to return it as JSON or YAML.

```
glooctl vhost get -o json

{"name":"vhost1","routes":[{"request_matcher":{"path_exact":"/bar","verbs":["GET","POST"]},"single_destination":{"upstream":{"name":"my-upstream"}}}],"metadata":{"resource_version":"226902","namespace":"gloo-system"}}
```

If you want to get details of a specific Virtual Host, you can use
`get` command with the name of the virtual host. It returns
the result as YAML, but you can use `output` flag to get JSON.

```
glooctl vhost get vhost1

metadata:
  namespace: gloo-system
  resource_version: "226902"
name: vhost1
routes:
- request_matcher:
    path_exact: /bar
    verbs:
    - GET
    - POST
  single_destination:
    upstream:
      name: my-upstream
```

### Updating Virtual Host
Similar to `create` command, `update` command takes the definition of
virtual host from a YAML file and replaces the existing virtual host with the
one from the file.

```
glooctl vhost update -f vhost2.yaml -o yaml

Virtual host updated
metadata:
  namespace: gloo-system
  resource_version: "228028"
name: vhost1
routes:
- request_matcher:
    path_exact: /bar
    verbs:
    - GET
    - POST
  single_destination:
    upstream:
      name: new-upstream
```

### Deleting Virtual Host
You can delete an existing virtual by giving the name of the
virtual host to be deleted to `delete` command.

```
glooctl vhost delete vhost1

Virtual host vhost1 deleted
```

## Reference
To learn more about Upstreams and Virtual Hosts please refer
to Gloo documentation.
