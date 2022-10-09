#!/bin/bash

kubectl create configmap ${app} --from-file=config.toml,local.toml
