# "Word of Wisdom" TCP Server with Proof-of-Work Based DDoS Protection

# Description:

"Word of Wisdom" TCP server with DDOS attack protection with Proof of Work 

# Types of requests
The system supports five types of requests, identified by a header:
+ 0 - Quit - signal to other side to close connection
+ 1 - RequestChallenge - from client to server - request new challenge from server
+ 2 - ResponseChallenge - from server to client - message with challenge for client
+ 3 - RequestResource - from client to server - message with solved challenge
+ 4 - ResponseResource - from server to client - message with useful info is solution is correct, or with error if not

# PoW:

As for PoW I decided to use [Hashcash](https://en.wikipedia.org/wiki/Hashcash) proof-or-work system.

PoW is implemented with a challenge-response protocol:

The client establishes a tcp connection with the server. The server starts to listening to client messages.

The client sends the RequestChallenge command to receive a challenge from server.

The client parses the message from the server, computes Hashcash, and submits the solution

The client requests a resource (in this case, a random quote) from the server.

# Choice of PoW

I had some options for Pow. For example: [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree), [Guided tour puzzle](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol)

After comparison, I chose Hashcash. Other algorithms have next disadvantages:
+ In Merkle tree server should do too much work to validate client's solution.
+ In guided tour puzzle client should regularly request server about next parts of guide, that complicates logic of protocol.

I have chosen hashcash for reasons:
1. Simple on implementation
2. More documentation on the web
3. Simpler for server validation

# How to get started
To start the project, ensure you have Docker and Docker Compose installed. Then clone the project and: 
```
make start
```
