# Tag checking and updating for alicloud on kubernetes

Version: v0.1.0

## Pre-Requirement
### RAM Permissions
```
{
    "Version": "1",
    "Statement": [
        {
            "Action": [
                "ecs:AddTag*"
            ],
            "Resource": "*",
            "Effect": "Allow"
        },
        {
            "Action": [
                "ecs:Describe*"
            ],
            "Resource": "*",
            "Effect": "Allow"
        },
        {
            "Action": [
                "vpc:Describe*"
            ],
            "Resource": "*",
            "Effect": "Allow"
        }
    ]
}
```
### Setup [Kube2ram](https://github.com/AliyunContainerService/kube2ram)

## Deploy
```
kustomize build kustomize/base | kubectl apply -f -
```

## Program Flags
```
Flags:
  -c, --cron string         cron scheduler
  -h, --help                help for updatek8stags
  -i, --instanceid string   filter by instance id
  -s, --pagesize int        alicloud api pagesize (default 10)

Global Flags:
      --config string     config file (default is $HOME/.ali-ecs-tag-update.yaml)
      --logfile string    log file
      --loglevel string   log level  [trace, debug, info, warn, error, fatal, panic] (default info) (default "info")
```      
