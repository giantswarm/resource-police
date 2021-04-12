# Resource Police

Informs via Slack about test clusters that may be deleted.

## Development

You can execute the job locally using the following command:

```nohighlight
CORTEX_USER_NAME=10755
CORTEX_PASSWORD=REDACTED
SLACK_WEBHOOK_ENDPOINT=REDACTED

go run . report
```
