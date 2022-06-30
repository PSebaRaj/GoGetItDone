# GoGetItDone
 Backend of a to-do app. Written in Go and utilizes Redis and PostgreSQL.

## Justifications for Technology Stack
### Redis
- Redis is used to cache tasks
	- Users will likely view a task's details frequently in a short span of time.

### PostgreSQL
- PostgreSQL is used for long-term storage of tasks
	- There will be many different types of tasks
	- Object inheritance for these types of tasks beyond the standard task
	- e.g. expiring tasks, priority tasks, overdue tasks

## Testing

### Clone Repository

### Postman
The Postman collection for testing the REST API functions can be found [here](https://www.getpostman.com/collections/40ab42d058be92ae4ef7):

To test:
1. Open Postman &rarr; Import &rarr; Link &rarr; Paste link from above &rarr; Import
2. Start backend, as described above

