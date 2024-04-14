# Tokeniser

A Golang + Postgres app that allows one to tokenise sensitive information. This personal project is aimed at teaching myself backend engineering, SRE, practical cryptography, and compliance. 

## How To Use

The current version of this app requires one to have the following environment variables, or a `.env` file in the same directory as the code:

```
PASSPHRASE=some-very-long-secret
POSTGRES_USER=my-database-user
POSTGRES_PASSWORD=my-secret-password
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=my-database-name
APP_PORT=8080
```

You can rename the `.env.example` file to `.env` and use it as well. However, please know that the default credentials aren't very secure, and you'll put yourself at risk.

The project uses a Makefile for setting up and running everything. To see all commands available, run the following in your terminal:

```console
$ make help
```

For either setup, ensure having `docker compose` and `make` installed. Before proceeding to the next step, you also need to ensure the `.env` file or the environment variables mentioned above have been set up. 

For now, I have not provided any way to check out the logs. However, if you choose to, you can view logs using the following commands:

```console
$ docker compose logs # For seeing logs from all services
$ docker compose logs app
$ docker compose logs db
$ docker compose logs flyway
```

In some future update, I'll setup a bash script that runs at startup, and will set up the environment variables by default in case the file and variables are missing. Along with it, I'd try to ease the process of seeing the logs. However, I'm yet to understand how to make this process easier so it might take some time.

### Local Development

Once the environment variables or `.env` file are set up, run the following command:

```console
$ make dev
```

### Running the application

For running the application, use the following command:
```console
$ make start
```

To get a fresh start, use the following command:
```console
$ make start-build
```

Note: This will not refresh your database. I have not provided the commands for dropping the DB for now, but might do it later.

## Workflow
User --> API --> Handler --> Database Handler -> Database

## Encryption

- A good read: [Practical Cryptography in Go](https://leanpub.com/gocrypto/read)
- The current encryption process begins with a passphrase. The passphrase, along with a salt, is used to derive key using [PBKDF2 (Password Based Key Derivation Function)](https://en.wikipedia.org/wiki/PBKDF2) that follows NIST specifications of 128 bits salt size with 210,000 iterations using SHA-512.
- The key derived from above method is used along with a nonce in AES-256 cipher with Galois-Counter mode to ensure security with performance. 
- Finally, the salt, and the nonce are prepended to the ciphertext. The final result is a byte slice containing salt, nonce, and ciphertext in that order.
- Decryption process follows similar pattern: derive salt, nonce and ciphertext. Use those values to obtain plaintext.

## Next steps:
- Use docker compose profiles for setting up environments
- Tests and Makefile are broken given there's a change in network of database. 
- Setup dev, testing, and prod all using make (prod done, use `make start` or `make start-build` for fresh build)
- Deployment of the API
- Improve performance. /all API takes 7 seconds for 15 cards. Setup concurrency.
- Better way for in-memory storage of sensitive information. Check this library: https://github.com/awnumar/memguard
- Implemetation of HTTPS
- Map all data to users
	- Keep one user until we work on user management
- Key management for encryption (Much later)

Note: Instead of writing my next steps here, I will begin using GitHub projects going ahead. It is accessible [here](https://github.com/users/arayofcode/projects/4).

## PCI-DSS Requirements:

| Number | Requirement | Status | Scope | Remarks |
|-|-|-|-|-|
| 3 | Setup data retention for 30 days | &#9744; | After 28th Feb | |
| 3 | Setup secure delete function for deletion after retention | &#9744; | After 28th Feb | |
| 3 | Setup encryption of keys that encrypt sensitive information | &#9744; | After 28th Feb | |
| 3 | Encrypt data | &#9745; | By 28th Feb | Using crypto library for encryption and decyption |
| 3 | Ensure sensitive data is not being logged anywhere intentionally or accidentally | &#9745; | By 28th Feb | No logging as of now. Still need to work on in-memory safety to avoid accidental dumping. Will check external libraries |
| 3 | Encrypt with strong cryptographic algorithms | &#9745; | By 28th Feb | Using AES-256 |
| 3 | Document Encryption and decryption process | &#9745; | Basic documentation by 28th Feb | Documented above |
| 3 | Mask sensitive information | &#9745; | By 28th Feb | Dashboard shows masked data. API part will be implemeted later |

## Deployment Strategy (Architecture-wise in GCP)

I'm not sure if I'd actually use GCP to deploy it (because it costs money :P) but here's the strategy:
- Two GCP Projects: 
  - **Infra:** For deployment of project. Don't see the need to deploy staging and prod in separate projects given separate instances could be deployed, and can configure two service accounts based on staging/ prod to access credentials from GSM.
  - **Sensitive:** Contains only sensitive information with minimum number of users. The DB containing sensitive information will be stored in this project. Access via service account created within this environment with minimum permissions. Any access to this project will be disabled otherwise, and no service inside this project will be allowed to access internet.
- The extended version of project would need two DBs: one that contains sensitive information, and the next one containing other data. For example, tokens or users. 
![alt text](docs/image-1.png)


## Think about later

- Possible Idempotency:
  - Request with same payload sent again
  - Same card details sent again. Should we generate new tokens in this case? 
- The current implementation in encryption generates a salt, uses it to derive key from a hash function. Then we generate a nonce, encrypt plaintext with the nonce, prepend nonce to the ciphertext. We then prepend salt to the ciphertext. 
  - Is it good to prepend salt to nonce + ciphertext, or keep salt with token in a separate table?
- Read underlying functions used in libraries
- Helpful for testing: https://www.paypalobjects.com/en_AU/vhelp/paypalmanager_help/credit_card_numbers.htm
- Change the "cipher" package name as well? It's clashing with another package that you're using
- Invalid expiry dates or credit card numbers would still work through dependency injection (using test cases)
  - Validate before encrypting and storing them
- Use `make` to create container, or run `make` commands within the containers?

## Useful References (aside from the documentation of tools):
- [One2N's SRE Bootcamp](https://playbook.one2n.in/sre-bootcamp): found most references there
- [Memguard](https://github.com/awnumar/memguard): For Memory-safe 
- [12 Factor Apps](https://12factor.net/)
- [Practical Cryptography in Go](https://leanpub.com/gocrypto/read)
- [PBKDF2 (Password Based Key Derivation Function)](https://en.wikipedia.org/wiki/PBKDF2)
- [Credit Card Numbers for Testing](https://www.paypalobjects.com/en_AU/vhelp/paypalmanager_help/credit_card_numbers.htm)
- [Docker Multistage Build](https://youtu.be/2QMoLyfIJx8?si=1EOyYpQytOaT_kas)
- [Creating a Postman Collection](https://youtu.be/NlrPjuXEaZ8?si=MAPg9KVYG5PogJmu)
- [Docker and Postgres with Persistent Data](https://youtu.be/G3gnMSyX-XM?si=ycVqtlaGYHmN7bjg)
- [Docker Networking](https://www.youtube.com/watch?v=OU6xOM0SE4o&ab_channel=HusseinNasser)

## Packages used during local development

```
make=4.4
golang=1.22
flyway-oss=10.10
postgres=16.2
docker=26.0.0
docker-compose=2.26.0
golangci-lint=1.57.2
hadolint
```