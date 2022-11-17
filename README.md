## ENV Vars

- store_password

## Config Vars

### http-url

Endpoint for upload to handle POST multipart-form request with "file" key containing file to be uploaded. JWT token is passed via authorization header containing relevant details about file. JWT is generated in InitTx process.

## TLS Cert

### Proposal

Contract stores hash of the self-signed x509 cert of the storage node to ensure clients are communicating with the actual storage node and avoid Man-in-the-middle attack.
