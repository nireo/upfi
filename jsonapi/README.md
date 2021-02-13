# JSON API documentation

**Routes needing a authorization token have a A before the method.**

### `POST /api/register`

Register a new account on the service.

Required fields:

- `username`
- `password`
- `master`

No optional fields

### `POST /api/login`

Login to your created user.

Required fields

- `username`
- `password`

No optional fields

### `POST A /api/file`

Upload a file for hosting. (This route doesn't take JSON but a multipart form)

Required fields

- `master`
- `file` (multipart file)

Optional fields

- `description`

### `GET A /api/files`

Get all of your files.

No required fields

No optional fields

### `GET A /api/file/:fileId`

Get more information about one file with a id

No required fields

No optional fields

### `GET A /api/me`

Get information about the user corresponding to the auth token.

No required fields

No optional fields

### `PATCH A /api/file/:fileId`

Update some fields on your file

No required fields

Optional fields

- `title`
- `description`

### `DELETE A /api/file/:fileId`

Deletes a file from hosting.

No required fields

No optional fields

### `DELETE A /api/me`

Deletes the account user who requests this.

No required fields

No optional fields

### `PATCH A /api/password`

Update users password.

Required fields

- `currentPassword`
- `newPassword`

No optional fields

### `POST A /api/download/:fileId`

Lets the user download a file from the server.

Required fields

- `master` (the file encryption master password)

No optional fields

### `PATCH A /api/settings`

Updates the settings of a user.

Required fields

- `username`

No optional fields
