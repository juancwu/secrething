# Konbini: your go to convenient store for your secrets

Konbini the backend service that stores secrets in a database and make them
accessible to your team, IC, and just yourself. The secrets are never saved in plaintext.
In fact, Konbini never knows the plaintext version of your secret.
The client is responsible of encrypting the secrets, which also means no key store in Konbini as well.
In combination with the CLI tool [Mi](https://github.com/juancwu/mi) the process is made even easier.

> Notice: The repository is undergoing a major do-over.
