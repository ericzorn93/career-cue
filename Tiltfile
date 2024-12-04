# Primary Docker Compose Imports
docker_compose(["./docker-compose.yml"])

# Accounts API
docker_build('career-cue/accounts-api', '.',
    dockerfile="./apps/services/accounts-api/Dockerfile",
    live_update = [
        sync('./apps/services/accounts-api', '/app'),
        run('air --build.cmd "go build -o /bin/accounts-api ./apps/services/accounts-api/cmd/server/main.go" --build.bin "/bin/accounts-api"'),
        restart_container()
    ]
)

# Accounts Worker
docker_build('career-cue/accounts-worker', '.',
    dockerfile="./apps/services/accounts-worker/Dockerfile",
    live_update = [
        sync('./apps/services/accounts-worker', '/app'),
        run('air --build.cmd "go build -o /bin/accounts-worker ./apps/services/accounts-worker/cmd/server/main.go" --build.bin "/bin/accounts-worker"'),
        restart_container()
    ]
)

# Inbound Webhooks API
docker_build('career-cue/inbound-webhooks-api', '.',
    dockerfile="./apps/services/inbound-webhooks-api/Dockerfile",
    live_update = [
        sync('./apps/services/inbound-webhooks-api', '/app'),
        run('air --build.cmd "go build -o /bin/inbound-webhooks-api ./apps/services/inbound-webhooks-api/cmd/server/main.go" --build.bin "/bin/inbound-webhooks-api"'),
        restart_container()
    ]
)