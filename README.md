# juanita-bot

## Docker Testing

### Build

`docker build -t juanita .`

### Run

`docker run --env-file ./.env --rm --name efs-viewer efs-viewer`

You have to create a `.env` file with the necessary environment variables set.

## Docker Deploy

`export REG_IP=<registry_ip>`
`docker build . -t ${REG_IP}:32000/juanita:latest`
`docker push ${REG_IP}:32000/juanita`
