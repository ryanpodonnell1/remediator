# User Story

As a Cloud Security Engineer I would like to ensure that I'm reducing the attack surface of public cloud assets by shutting down wide open sensitive ports

## Proactively Preventing this type of issue

In order to reduct attack surface it is recommended to disable ingress ports on known sensitive ports such as 22(SSH), 3389(RDP), 3306(MySQL)

### Prevent:

One method to prevent this scenario is to institute a Pipeline mechanism that will check IaaC output as a preventative measure. Such tooling that can be employed is OPA(Open Policy Agent) [Terraform Plan Example](https://gist.github.com/ryanpodonnell1/3da9805733ce7dcce71ee5e0622fb1cc) and Checkov. These types of checks should be insituted in PR builds to give developer feedback as quick as possible and prevent git merges of known bad configuration

Checkov is also available in IDEs such as VSCODE that can catch this type of vulnerability before even being committed to source on the developers workstation. Other tooling that is not available in IDEs can be employed via pre-commit hooks to *validate* the code prior to commit, such as executing an OPA policy against a local tf plan in a test environment

More advanced examples may be to provision approved security groups through firewall manager and prevent new security groups being provisioned anywhere else. This would provide central visibility and central source for networking traffic.

### Detect:

Not all configuration is done via IaaC (Terraform, Cloudformation, etc) and can be manually done through the console. Preventative checks in CI/CD Pipelines may be limited when this occurs. Ad hoc checks such as config rules, lambdas, cloudcustodian should occur on a regular basis to ensure these type of configurations aren't present in an environment in the event they have circumvented the CI checks.

### Report:

Vulnerabilities such as these should be reported on a regular basis to establish whether the remediation/preventative controls are effective.

## Authenticating to AWS

It's recommended to use `export AWS_PROFILE=<profile_name>`, there are also variable placeholders in the Makefile if you wish to use the terraform/cli commands

## Spinning up the infrastructure with Provided Terraform  

A terraform configuration has been provided to make deployment/testing/teardown easier and is not required for the remediation code to work. Ensure that you are running **terraform v14** if you wish to use `make tf_apply`


### Manual Steps

If you don't wish to use the terraform provided, configure 1 or more AWS Security Group Ingress rules that match `0.0.0.0/0` and any port `[22,3389,3306]`. They may also fall within ranges such as `0-1024`


## CLI Tool

### Build

`make build` will place the compiled program in `bin/`

### Usage

The cli tool is called remediator and can be used several ways

| Command                                | Description                                              |
| -------------------------------------- | -------------------------------------------------------- |
| `remediator detect`                    | prints out all securitygroups that have                  |
| `remediator remediate`                 | performs a dry-run of removing resources                 |
| `remediator remediate --dry-run=false` | performs an active run to remove violating ingress rules |


Makefile excerpt for usage examples:
```makefile
detect: build
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd bin && \
	./remediator detect

remediate_dry: build
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd bin && \
	./remediator remediate --dry-run=true

remediate_active: build
	export AWS_PROFILE=$(DEFAULT_AWS_PROFILE) && \
	cd bin && \
	./remediator remediate --dry-run=false
```


## TODOs

* add multiple outputs for better parsing/reporting (i.e. JSON, CSV, etc)
* add --auto-approve/manual approval mechanism (similar to terraform)
* add ability to consume pre-planned files for approval workflow, this ensures that what is approved is what is actually removed 
* add notification mechanism such as webhooks, email, etc