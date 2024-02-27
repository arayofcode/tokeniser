## Workflow
User --> API --> Handler --> Database

## Encryption

- [Practical Cryptography in Go](https://leanpub.com/gocrypto/read)
- 



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

| Number | Requirement | Status | Scope |
|-|-|-|-|
| 3 | Setup data retention for 30 days | &#9744; | After 28th Feb |
| 3 | Setup secure delete function for deletion after retention | &#9744; | After 28th Feb |
| 3 | Setup encryption of keys that encrypt sensitive information | &#9744; | After 28th Feb |
| 3 | Encrypt data | &#9744; | By 28th Feb |
| 3 | Ensure sensitive data is not being logged anywhere intentionally or accidentally | &#9744; | By 28th Feb |
| 3 | Encrypt with strong cryptographic algorithms (AES?) | &#9744; | By 28th Feb |
| 3 | Document Encryption and decryption process | &#9744; | Basic doc by 28th Feb |
| 3 | Mask sensitive information | &#9744; | By 28th Feb |

## Think about later

- The current implementation in encryption generates a salt, uses it to derive key from a hash function. Then we generate a nonce, encrypt plaintext with the nonce, prepend nonce to the ciphertext. We then prepend salt to the ciphertext. 
  - Is it good to prepend salt to nonce + ciphertext, or keep salt with token in the mapping?