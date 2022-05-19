#!/bin/bash

kubectl create configmap ${app}-config --from-file=config.toml,local.toml
