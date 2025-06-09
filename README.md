# Discord Bot Backend
Open-sourced code for the Typical Developers Discord bot API.

### Environment
Refer to the `.env.example` for environmental variables.

## Developing
### Prerequisites
- [Golang 1.23+](https://go.dev/)
- Docker
- [Taskfile](https://taskfile.dev/)
Run `go install` to install dependencies.

### Deploying API
```
task dev:api
```
Optionally, if you want to run the API without building docs and generating SQL schemas:
```
task dev:compile-api-only
```

### Developing Scheduled Tasks
```
task dev:tasks
```

### SQL Migrations
Migration are handled by [migrate](https://github.com/golang-migrate/migrate). Migrations should be created in `internal/db/migrations` and can be ran with `task migrate:up` (to update) or `task migrate:down` (to rollback). Queries should be tested before being pushed to production.

## Licensing
All code for the bot is licensed under the [GNU General Public License v3.0](https://github.com/typical-developers/main-discord-bot/blob/main/LICENSE) license. Please refer to the LICENSE file for more information regarding rights and limitations.

TL;DR: You are allowed to do whatever with the code (modify, sell, redistribute, etc) as long as you allow others to do the same with yours.

## Resources
- [Typical Developers Discord Server](https://discord.gg/typical)
- [Discord Bot](https://github.com/typical-developers/main-discord-bot)