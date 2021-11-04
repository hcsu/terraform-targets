# Terraform Targets

Operate (plan or apply) multiple targets without `-target` subcommand.

```sh
#!/usr/bin/env bash

ARGS=("${@:2}")
targets=${ARGS[*]/#/-target }
terraform "$1" $targets
```

Usage:
```sh
tt plan 'module.foo.ep_association_s3[0]' 'module.bar.frontend_https[0]'
```
