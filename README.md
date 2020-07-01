# Resource Police

Warning test cluster owners about excess resource consuption

## Configuration

The tool needs an installation config file with one entry per installation. Example:

```yaml
- name: myinst
  apiEndpoint: https://API-ENDPOINT
  credentials:
    password: redacted
    user: redacted
```

## Development

You can execute the job locally using the following command:

```nohighlight
go run . report \
  --installations.config.file ./config.yaml \
  --slack.webhook.endpoint https://hooks.slack.com/services/redacted
```
