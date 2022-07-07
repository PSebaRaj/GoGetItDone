# GoGetItDone
[![Latest Release](https://img.shields.io/github/release/psebaraj/gogetitdone.svg?style=for-the-badge)](https://github.com/psebaraj/gogetitdone/releases)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE)
[![Build Status](https://img.shields.io/github/workflow/status/psebaraj/gogetitdone/Go?style=for-the-badge)](https://github.com/PSebaRaj/GoGetItDone/actions/workflows/go.yml)
[![Go ReportCard](https://goreportcard.com/badge/github.com/psebaraj/gogetitdone?style=for-the-badge)](https://goreportcard.com/report/psebaraj/gogetitdone)
[![Lines of Code](https://img.shields.io/tokei/lines/github/psebaraj/gogetitdone?style=for-the-badge)](https://github.com/psebaraj/gogetitdone/actions)

Backend of a to-do app. Written in Go and utilizes Redis and PostgreSQL.

## To-Do:
- [x] Error handling on concurrent processes
	- Namely UpdateExpiringTask (updates DB and JSON res)
- [x] Set/change cache expiration time for REDIS
- [ ] Use Prepare (Postgres) to cache PSQL statements for speed (v1.1)
	- UpdateExpiringTask
	- [GORM Doc 1](https://gorm.io/docs/performance.html)
	- [GORM Doc 2](https://gorm.io/docs/v2_release_note.html#Prepared-Statement-Mode)
- [ ] Refactor tasks with interfaces (v1.1)

## Features
- Getting an expiring task will return the expiration time in the local time zone of the origin of the request
	- Give expiry time for expiring tasks in GMT for consistency
- Getting a user returns their details and all of their tasks
	- (Regular) tasks, expiring tasks, priority tasks, etc.

## Justifications for Technology Stack
### Redis
- Redis is used to cache tasks
	- Users will likely view a task's details frequently in a short span of time.
	- Note that only the tasks are cached, not the user data

### PostgreSQL
- PostgreSQL is used for long-term storage of tasks
	- There will be many different types of tasks
	- Object inheritance for these types of tasks beyond the standard task
	- e.g. expiring tasks, priority tasks, overdue tasks
- User will have a one-to-many relationship with tasks

## Testing
### Redis
- Install Redis, if not already installed
- To install:
	- `brew install redis`
	- `brew services start redis`
	- `redis-server /usr/local/etc/redis.conf`
- Check if Redis is running and its port:
	- `ps -ef | grep redis`
	- ![RedisCheck](./pictures/CheckRedisRunning.png)

### PostgreSQL
- Install Postgres, if not already installed
	- `brew install postgresql`
	- `initdb /usr/local/var/postgres/`
- Start Postgres
	- `sudo psql -U my_macosx_username postgres`
	- `brew services start postgresql`
	- `ALTER USER my_macosx_username PASSWORD 'new_password';`
	- `CREATE DATABASE todo_list;`

### Start Backend
- To run the backend of the application, first clone the repository:
	- `git clone https://github.com/psebaraj/gogetitdone.git`
- Navigate to the GoGetItDone directory
- Create .env file and entire appropriate credentials for:
	- `DB_DIALECT, DB_HOST, DB_PORT, DB_USER, DB_NAME, DB_PASSWORD`
	- `REDIS_HOST, REDIS_PORT, REDIS_PASSWORD`
- Build and run the application:
	- `go run main.go`

### Postman
The Postman collection for testing the REST API functions can be found [here](https://www.getpostman.com/collections/40ab42d058be92ae4ef7)
