# Server

Contains go-do-list server logics.

### Local Testing
To be able to run tests locally, there are some prerequisites data that should be provided and exported:

```bash
$ export DATABASE_CONNECTION="user=postgres password=password host=localhost port=5432 dbname=postgres connect_timeout=5 sslmode=false"
```