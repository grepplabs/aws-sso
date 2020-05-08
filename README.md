# aws-sso
AWS Single Sign-On utilities

### Prerequisites
- [aws cli version 2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)

## Install binary release

1. Download the latest release

   Linux

        curl -Ls https://github.com/grepplabs/aws-sso/releases/download/v0.0.1/aws-sso-v0.0.1-linux-amd64.tar.gz | tar xz

   macOS

        curl -Ls https://github.com/grepplabs/aws-sso/releases/download/v0.0.1/aws-sso-v0.0.1-darwin-amd64.tar.gz | tar xz

   windows

        curl -Ls https://github.com/grepplabs/aws-sso/releases/download/v0.0.1/aws-sso-v0.0.1-windows-amd64.tar.gz | tar xz

2. Move the binary in to your PATH.

    ```
    sudo mv ./aws-sso /usr/local/bin/aws-sso
    ```

## Build binary

    make clean build


### Usage
#### Initial setup

1. AWS config example

    ~/.aws/config

   ```
    [profile development]
    sso_start_url = https://d-4711.awsapps.com/start
    sso_region = eu-central-1
    sso_account_id = 000000000000
    sso_role_name = aws-developer
    region = eu-central-1
    output = json
    cli_pager=
    ```

2. SSO login

   ```
    $ aws sso login --profile development
   ```

### aws-sso credentials commands

#### export

   ```
   $ aws-sso credentials export --profile development

   export AWS_ACCESS_KEY_ID="your_access_key_id"
   export AWS_SECRET_ACCESS_KEY="your_secret_access_key"
   export AWS_SESSION_TOKEN="your_session_token"


   $ eval $(aws-sso credentials export --profile development)
   ```

#### refresh

   ```
   $ touch ~/.aws/credentials
   $ aws-sso credentials refresh --profile development
   $ cat ~/.aws/credentials

   [development]
   aws_access_key_id = your_access_key_id
   aws_secret_access_key = your_secret_access_key
   aws_session_token = your_session_token
