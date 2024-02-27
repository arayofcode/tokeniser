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

### Scope 
- Systems performing encryption and/or decryption of cardholder data, and systems performing key management functions,
- Encrypted cardholder data that is not isolated from the encryption and decryption and key management processes,
- Encrypted cardholder data that is present on a system or media that also contains the decryption key,
- Encrypted cardholder data that is present in the same environment as the decryption key,
- Encrypted cardholder data that is accessible to an entity that also has access to the decryption key

### Requirement 1: Install and Maintain Network Security Controls 

- Out of Scope for now

### Requirement 2: Apply Secure Configurations to All System Components

- Out of Scope for now

### Requirement 3: Protect Stored Account Data

- If account data is present in non-persistent memory, encryption of account data is not required
- Remove data as soon as the transaction completes
- Ensure the process and mechanism for protecting stored account data is well defined and documented.
- Define Roles and Responsibilities for performing activities
- Minimise storage of account data 
- 3.5 out of scope
- Where cryptography is used to protect stored account data, key management processes and procedures covering all aspects of the key lifecycle are
defined and implemented.

#### Minimise Storage of Account Data

- Implementation of data retention and disposal policies
- The policies should cover location of stored account data, as well as data stored before authorisation
- Limit data storage amount and retention time to what regulation allows
- Secure delete, or make account data irrecoverable after retention period is over
- Process to verify the above point at least every three months
- Data allowed to retain even after retention period is over: PAN, name, expiration date and service code
- Consider processes and users with access to sensitive data
- Don't overlook backup, archives, removable data storage devices, paper-based media and audio recordings

#### Sensitive authentication data (SAD) is not stored after authorization.

- Remove the data safely as early as possible to avoid potential data leakage during or after the transaction
- The full content should never be retained after completion of authorisation process
- Review: Transaction data, logs, history file, trace file, db schema, db content (on-premise and cloud data store), memory/ crash or dump files
- CVV, PIN and PIN Block are not retained after authorisation

#### Sensitive Authentication Data stored during authorisation is encrypted using Strong Cryptography

- The authorization process is completed as soon
as the response to an authorization request
response—that is, an approval or decline—is
received.
- Encryption Key for SAD and PAN should be different
- Storage of sensitive auth data must be limited to what is needed for business need, and encrypted using strong cryptography
- Mask PAN details so that only people with right business need can see more than BIN and last four digits
- Prevent copying of the PAN details
- PAN is secured wherever it is stored

#### Cryptographic Keys are Secured

- Access to keys is minimised to minimum number of people necessary (people creating, altering, rotating,
distributing, or otherwise maintaining encryption
keys)
- Key-encrypting keys are at least as strong as the ones they secure
- Key-encrypting keys are stored separately from the data-encrypting keys
- Keys are stored in fewest possible locations and forms
- Maintain documentation about:
  - Details of all algorithms, protocols, and keys
used for the protection of stored account data,
including key strength and expiry date.
  - Same key should not be used in prod and staging
  - HSMs, Key management systems and secure cryptographic devices
- Secrets and private keys must be kept in one of the following form
  - Encrypted with a key that is at least as strong, and stored separately from the data encrypting key
  - Within secure cryptographic device/ HSM, etc
  - At least two full key-length components or key shares (partial keys)
- When using HSM, HSM interaction channel must be secured too
- 

### Requirement 4: Protect Cardholder data with Strong Cryptography During Transmission over Open, Public Networks

- Out of Scope for now

### Requirement 5: Protect All Networks and Systems From Malicious Software

- Out of Scope for now

### Requirement 6: Develop and Maintain Secure Systems and Software

- Out of Scope for now

### Requirement 7: Restrict Access to System Components and Cardholder data by Business Need to Know

- Out of Scope for now

### Requirement 8: Identify Users and Authenticate Access to System Components

- Out of Scope for now

### Requirement 9: Restrict Physical Access to Cardholder Data

- Out of Scope for now

### Requirement 10: Log and Monitor All Access to System Components and Cardholder Data

- Out of Scope for now

### Requirement 11: Test Security of Systems and Networks Regularly

- Out of Scope for now

### Requirement 12: Support InfoSec with Organizational Policies and Programs

- Out of Scope for now