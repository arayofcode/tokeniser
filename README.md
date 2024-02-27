## Workflow
User --> API --> Handler --> Database

## Possible Idempotency:
- Request with same payload sent again
- Same card details sent again. Should we generate new tokens in this case? 

## Next steps:
- Encryption of secrets received
- Secure token generation instead of UUIDs
- Implemetation of HTTPS
- Map all secrets to user
	- Keep one user until we work on user management
- Key management for encryption

## PCI-DSS Requirements:

### Scope 
- Systems performing encryption and/or decryption of cardholder data, and systems performing key management functions,
- Encrypted cardholder data that is not isolated from the encryption and decryption and key management processes,
- Encrypted cardholder data that is present on a system or media that also contains the decryption key,
- Encrypted cardholder data that is present in the same environment as the decryption key,
- Encrypted cardholder data that is accessible to an entity that also has access to the decryption key

### Approaches to Implement

#### Defined Approach

- Meet the stated requirements
- 