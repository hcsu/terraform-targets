# Terraform Targets

Operate (plan or apply) multiple targets without `-target` subcommand.

## Installation

```sh
go build -o tt
# Or install to $GOPATH/bin
go install
```

## Usage

```sh
# Extract and list all targets from plan
tt -t

# Extract and list module-level targets (deduplicated)
tt -m

# Apply terraform command with multiple targets
tt plan 'module.foo.ep_association_s3[0]' 'module.bar.frontend_https[0]'
tt apply 'module.foo.ep_association_s3[0]'
```

## Examples

For a plan output containing:
```
# module.web_app.foo.bar[0] will be created
# module.web_app.baz.qux[0] will be created
# aws_s3_bucket.example will be updated
```

`tt -t` outputs:
```
'module.web_app.foo.bar[0]' \
'module.web_app.baz.qux[0]' \
'aws_s3_bucket.example'
```

`tt -m` outputs:
```
'module.web_app' \
'aws_s3_bucket.example'
```
