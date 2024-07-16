This system will allow users to upload documents, which are then hashed and stored on the blockchain. Other users can verify the authenticity of these documents by comparing their hashes.

## Problem Statement
In many industries, verifying the authenticity of documents is crucial. For example, in legal, educational, or financial sectors, ensuring that a document has not been tampered with is vital. Blockchain technology can provide a tamper-proof solution by storing document hashes on a decentralized ledger.

## Solution Overview
We will create a document verification system with a Go backend and an HTML/CSS frontend. The system will allow users to:

- Upload documents.
- Store the document hashes on the blockchain.
- Verify the authenticity of uploaded documents.
## Explanation
1. **Backend (Go):**

- Defines Document and Blockchain structures.
- Implements hashing and storing document hashes on the blockchain.
- Implements RPC methods to add documents and get the document history.
- Provides an HTTP server to interact with the blockchain via a web interface.
- Implements client functions to add documents and retrieve the document history via RPC.
2. **Frontend (HTML & CSS):**

- Provides a form to upload new documents.
- Displays the current document history in a list format.

This setup allows you to run a blockchain-based document verification system that ensures the integrity and authenticity of documents by storing their hashes on an immutable blockchain ledger.