## Next steps
- Use [Postgres ts vector](https://www.postgresql.org/docs/current/datatype-textsearch.html) data type to provide a full text search for users.
- Implement login and use JWT to interact with API
- Refactor code to use connection pool for database
- Seems we have duplicated email on csv file, maybe will be better refactor code to have unique INDEX on database and `DO NOTHING` for duplicate entries during import data
- ⚠️ Troubleshooting the problem with consumer (always getting 0 messages)
- Finish [PR related load balancer](https://github.com/rdiego26/golang-users-api/pull/1)