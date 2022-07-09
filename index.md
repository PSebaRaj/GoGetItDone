## GoGetItDone - Github Pages
[![Latest Release](https://img.shields.io/github/release/psebaraj/gogetitdone.svg?style=for-the-badge)](https://github.com/psebaraj/gogetitdone/releases)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE)
[![Build Status](https://img.shields.io/github/workflow/status/psebaraj/gogetitdone/Go?style=for-the-badge)](https://github.com/PSebaRaj/GoGetItDone/actions/workflows/go.yml)
[![Go ReportCard](https://goreportcard.com/badge/github.com/psebaraj/gogetitdone?style=for-the-badge)](https://goreportcard.com/report/psebaraj/gogetitdone)

Backend of a to-do app. Created because I don't want to pay for Todoist Premium (and there are some missing features surrounding collaboration that I really want).

Written in Go and utilizes Redis for caching and PostgreSQL for long-term storage. Currently on v1.0.0, which implements users (people) and three different types of tasks. Currently in development, v1.1.0 will include refactored tasks, as well as two new features (projects, which are groups of tasks, and groups, which are collections of users that can share tasks). More information about v1.1.0 can be found [here](https://github.com/PSebaRaj/GoGetItDone/blob/main/README.md#to-do).

As of right now, only the backend (server and cache) is complete. A friend has agreed to create the client, and to make it as painless as possible, I have documented all of the endpoints with Swagger.

### Languages, Technologies, and Frameworks
#### Server
- Go
- Gorilla Mux
- PostgreSQL

#### Cache
- Redis
