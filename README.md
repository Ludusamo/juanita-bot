# juanita-bot

## Docker Testing

### Build

`docker build -t juanita .`

### Run

`docker run --env-file ./.env --rm --name juanita juanita`

You have to create a `.env` file with the necessary environment variables set.
