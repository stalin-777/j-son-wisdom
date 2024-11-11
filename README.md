# J-son-wisdom

Word of Wisdom TCP Server and Client with Proof of Work (PoW) that responds with a Jason Statham "quote". Implemented as part of a test task

## Task Description

Design and implement "Word of Wisdom" tcp server.

- TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from "word of wisdom" book **or any other collection of the quotes**.
- Docker file should be provided both for the server and for the client that solves the POW challenge

## Architecture and Implementation Overview

### Key Components

- **TCP Server**: Listens on a port, generates a PoW challenge, sends it to the client, and verifies the received solution. If the solution is correct, the server sends a random quote to the client.
  
- **Proof of Work (PoW)**: Used to protect against DDoS attacks. Before providing a quote, the server sends a computational challenge requiring the client to find a `nonce` (random value) such that the hash of the challenge combined with the nonce has a specified number of leading zeroes (e.g., 5).
  
- **Quote Collection**: After a successful PoW verification, the server selects a random quote from a list and sends it to the client.

- **Client**: Connects to the server, solves the PoW challenge, and sends the solution to the server. Upon successful verification, the client receives a quote.

## System Design and Implementation

### Server Workflow

1. **PoW Generation**: The server generates a random string as a "challenge" and sends it to the client.
2. **PoW Solution**: The client computes a `nonce` that, when combined with the challenge, produces a hash starting with the specified number of zeroes.
3. **PoW Verification**: The server verifies the hash of the solution. If the result is valid, the client receives a random quote.
4. **Quote Sending**: Upon successful PoW verification, the server sends a random quote to the client.

### Hashing algorithm

In this implementation, you can use 2 hash function options: sha1 and sha256. sha1 was chosen because it is used in the classic implementation of Hashcash, and sha256 is a more modern algorithm used in Bitcoin, among other things. If necessary, the list can be expanded

---

## Installation and Running

### 1. Running the Server

To simply start a project, use the command:

```make start```

Other commands can be found in the Makefile.
