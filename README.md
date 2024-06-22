# Terraform Targets

Operate (plan or apply) multiple targets without `-target` subcommand.

```sh
#!/usr/bin/env bash

# Check if the -t flag is provided
if [[ "$1" == "-t" ]]; then
  terraform plan | ag 'will be' | awk -F'# ' '{print $2}' | awk -F' will be' '{print $1}' | awk '{if(NR>1)print prev " \\"; prev=$0} END {print prev}'
  exit 0
fi

# Continue with the original script if the -t flag is not provided
ARGS=("${@:2}")
targets=${ARGS[*]/#/-target }
terraform "$1" $targets
```

Usage:
```sh
tt plan 'module.foo.ep_association_s3[0]' 'module.bar.frontend_https[0]'
tt -t
```
