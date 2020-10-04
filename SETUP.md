# Setup

Instructions on how to set up the benchmarks.

## Prerequisites

- Docker version >= 19.03.13
- docker-compose version >= 1.27.4 (usually included in the Docker installation)
- go version >= 1.14.2

## Local variables

The following variables should be set (use this in your `.bashrc`, `.zshrc` etc)

```bash
export GOPRIVATE="github.com/mmathys"
export AWS_DEFAULT_OUTPUT="text"
export DOCKER_BUILDKIT=1
```

## GitHub

1. Make sure that you have write access to `github.com/mmathys/acfts`
1. Make sure that your Git remote is over SSH, not HTTPS
1. Make sure that your public key in `~/.ssh/.id_rsa.pub` is added to your GitHub account

## AWS

1. Create an AWS account
1. Install the AWS CLI and make sure the binary is in your `PATH`
1. Configure the CLI for your account with `aws configure`
1. Generate a key pair for the selected AWS region and store it on your machine (I used `us-west-1`)
1. Create a security group
    1. Name it `all-open`
    1. Create inbound and outbound rules: all traffic should be allowed.
1. Create two separate EC2 launch configurations for validators and clients and save the launch configuration IDs.
    1. The image should be an Amazon Linux 2 AMI
    1. I used EC2 compute instances (beginning with `c`) but you can basically freely chose any type.
    1. Make sure that your generated keypair is added to both launch configurations.
    1. the security group should be the one we created before: `all-open`
1. Make sure that your account is able to spin up enough CPUs for the selected region (you may want to create a support
request to increase your limit)

## Configure `test/launch_scripts/awscli/launch.go`

1. Replace the strings in `ValidatorLaunchTemplate` and `AgentsLaunchTemplate` to your according launch configuration
IDs.
1. Modify `SSHKey`: insert the path of your own AWS keypair
1. Modify `SedExe`: change it to `sed`, because you're probably using Linux. If you're using MacOS, install `gnu-sed`
from Homebrew and set `SedExe` to `gsed`.