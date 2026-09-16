#!/bin/bash -x

___INSTANCE_NAME___="{{ container.name }}"

{{ config_get("user.kind.cloud-init-launch.sh", properties.default) }}
