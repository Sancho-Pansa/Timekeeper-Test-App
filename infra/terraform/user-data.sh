#! /bin/bash

set -euxi pipefall

dnf install -y docker awscli-2

systemctl enable --now docker
systemctl enable --now amazon-ssm-agent