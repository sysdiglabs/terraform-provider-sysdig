OANC Script
===========

This script allows you to retrieve the Notification Channels from Sysdig Secure that must be created via UI because
they enforce OAuth.

Steps to retrieve the resources:

1. Create the Notification Channel via UI (eg. Slack channel)

2. Modify the docker-compose.yaml file with your own SDC_TOKEN.

3. Execute `docker-compose run --rm oanc`.

This will print the resources in your CLI.

Just copy and paste the results in your Terraform file, and, if needed, remove the channels created via UI.