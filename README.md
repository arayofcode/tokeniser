## Workflow
User --> API --> Handler --> Database

## Possible Idempotency:
- Request with same payload sent again
- Same card details sent again. Should we generate new tokens in this case? 

## Next steps:
- Add test for handlers
- Data Validation
- Encryption of secrets received
- Secure token generation instead of UUIDs
- Implemetation of HTTPS
- Map all secrets to user
	- Keep one user until we work on user management
- Key management for encryption