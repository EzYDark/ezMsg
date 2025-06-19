# Final Architectural Summary: ezMsg

## I. Core Principles

* **Absolute Zero-Knowledge Server:** The server has zero knowledge of any user secrets. It only stores public keys and non-sensitive metadata, acting as a secure message relay and directory.
* **Locally-Protected, Device-Bound Keys:** A user's identity is based on a `privateKey` that is generated and stored on their device. This key is always encrypted at rest, protected by a strong local password set by the user for that specific device.
* **Self-Custody via Recovery Phrase:** The one and only method for account recovery after losing all devices is a user-held, 24-word recovery phrase. The user has full control and full responsibility.
* **Secure Device Linking & MFA:** Adding new devices or performing a recovery is a cryptographic ceremony, further protected by an optional second factor (TOTP) for high-stakes actions.
* **Efficient, Multi-Device E2EE:** End-to-end encryption works seamlessly across all devices and in group chats, using an efficient message delivery model to save bandwidth.

## II. Data Storage Model

* **Server-Side (Decoupled Model):**
    * **`Users` Table:** Contains `UserID`, `Username`, an encrypted `TOTPSecret`, and `IsTOTPEnabled`.
    * **`Devices` Table:** Contains `DeviceID`, `UserID`, and `PublicKey`.
    * **`Messages` Table:** Stores the `EncryptedContent` ("lockbox") and a list of `MessageKeyIDs`.
    * **`MessageKeys` Table:** Stores individual `EncryptedKeyPayload`s ("safes").
    * **`RecoveryCheckpoints` Table:** Stores the `encryptedHistoryBlob` for each user after a key migration.

* **Client-Side:**
    * A private file containing the `encryptedPrivateKey`, protected by a key derived from the local password via **Argon2id**.
    * An **encrypted SQLite database** (e.g., using SQLCipher). This file is unlocked with a key derived from the local password and is read from efficiently without loading the entire file into memory.

## III. The Authentication & Device Management Flow

#### A. First-Time User Registration & Recovery Setup
1.  **Key Generation & Recovery Phrase:** The app generates a new key pair and the corresponding 24-word recovery phrase, forcing the user through a multi-step verification process (including a comprehension check and word confirmation).
2.  **Create Local Password:** After verification, the user creates a strong local password (enforced by UI rules) to secure the device.
3.  **Local Encryption & Server Update:** The app uses **Argon2id** on the password to encrypt the `privateKey` and stores it locally. The `publicKey` is sent to the server.
4.  **Prompt for TOTP:** After registration, the app should strongly encourage the user to set up TOTP as a second factor for enhanced security.

#### B. Subsequent App Launch & Session Renewal
1.  **Unlock Private Key:** The user opens the app and enters their local password to decrypt the `privateKey` into memory and unlock the local database.
2.  **Challenge-Response:** To get a new access token, the client signs a server `challenge` with the in-memory `privateKey`.
3.  **Token Issuance:** The server verifies the signature and issues a new, **short-lived `accessToken`** with a configurable expiration (e.g., default 60 minutes).

#### C. Adding a New Device (The "Linking" Ceremony)
1.  **Initiation:** A new device generates a key pair and displays its `publicKey` as a QR code and a manual text phrase as a fallback.
2.  **Approval & MFA Check:** An active device scans the code or uses the phrase, gets user approval via their local password, and then performs a **TOTP check** with the server (if the user has TOTP enabled).
3.  **Server Update:** After the MFA check succeeds, the active device signs an approval message. The server verifies this signature and adds the new device's `publicKey` to the user's account.

#### D. Removing a Linked Device
1.  An active device sends a **signed command** to the server to remove another one of the user's devices.
2.  The server verifies the signature (and can optionally require a TOTP check).
3.  If valid, the server deletes the target device's `PublicKey` and also performs a surgical deletion on the `MessageKeys` table, removing all historical keys associated with the removed device.
4.  The server then broadcasts a "device removed" event to all of the user's chat partners.

## IV. The E2EE Messaging Flow

1.  **Sending (Atomic Operation):** The client assembles **one single, atomic package** containing the encrypted message content ("lockbox") and the list of encrypted `MessageKeys` ("safes") and sends it to the server in one request. The server parses this package and stores the data in its decoupled `Messages` and `MessageKeys` tables.
2.  **Receiving (Efficient Delivery):** The server sends a tailored package to each recipient containing the encrypted message and **only the single encrypted `MessageKey`** intended for that device.

## V. Dynamic Chat Management

* **When a Member Adds/Removes a Device or Recovers their Account:** The server broadcasts a system event (e.g., "device added," "key changed"). Receiving clients must listen for this event and automatically update their local record of the user's public keys. The UI must display a security notification for any key change.
* **When a New User is Added to a Chat:** The server sends a "new user added" event. The new user **does not** get access to the message history from before they joined.

## VI. Account Recovery & Key Migration

This flow handles the "lost all devices" scenario.

1.  **Initiation:** A user enters their **old 24-word recovery phrase**.
2.  **MFA Check:** Before proceeding, the server requires a **TOTP code** (if the user had it enabled) to authorize the recovery attempt.
3.  **Generate New Identity:** After the MFA check passes, the app regenerates the `old_private_key`, then immediately generates a `new_private_key` and a **new 24-word recovery phrase**, forcing the user through the backup and verification process.
4.  **Update Server Identity:** The client informs the server of its `new_public_key`. The server invalidates all old device keys and **broadcasts a "security key changed" event** to all chat partners.
5.  **Create New History Checkpoint:** The client performs a one-time migration: it decrypts the entire chat history locally, creates a single `encryptedHistoryBlob` from the plaintext, encrypts it with the `new_private_key`, and uploads it to the server as the **master recovery checkpoint**.
6.  **Invalidate Old Message Keys:** The client sends a final, signed command to the server to delete all `MessageKeys` associated with the old devices.
7.  **Finalize:** The `old_private_key` is securely wiped from memory. The old recovery phrase is now useless for both future authentication and for decrypting past messages.

## VII. Critical Security & Implementation Considerations

* **Rate Limiting:** All server actions (login attempts, account creation, etc.) must be strictly rate-limited.
* **Challenge Expiration:** Authentication challenges must be single-use and expire quickly.
* **Chat Permissions:** A permission model must be enforced on the server for who can add or remove members from a group chat.
* **Forward Secrecy:** For future enhancement, implement a Double Ratchet algorithm to provide Forward Secrecy, protecting past messages even if a private key is compromised.
* **Metadata Privacy:** Be aware that while message *content* is private, your server can still see *metadata* (who is talking to whom, when, etc.).