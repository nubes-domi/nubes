# Nubes storage

To meet the requirement for a simple backup, all the apps must share a single file-system backed datastore.

Whenever a file-system like path is used (eg: /my/pictures/2022), the root directory is to be intended as the root of the datastore.

Broadly, there's 3 categories of accesses that we should support:

- Relational data
- Object (S3 like) storage
- File (open/read/write/close system calls) storage

## Relational data

All relational data MUST be stored in SQLite databases.

Each application MAY create up to one SQLite database. The database file must be stored at `/sql/<application codename>.sqlite`

Applications SHOULD NOT access each other database files, using APIs instead to communicate.

Data that fits the relationa model:
- Comments on a picture
- Auditing for login attempts
- Star rating on a song

## Object storage

In a lot of scenarios, data files can be considered as blobs of data that once written, can only be fully replaced or deleted. Updating a portion of the object is not supported.

Data that should be stored as objects, anything not relational:
- documents
- pictures
- songs
- spreadsheets

## File storage

In some small, specific scenarios replacing an entire file for a small change could be an overkill. For example, changing the metadata on a 20 Gb video file.

For these cases, the regular open/seek/write/close system calls should be used.

Additionally, users might expect to browse their data using the regular operating system file browser and/or command line tools. This is also required to perform backups and restore.

## File storage over the network

The entire datastore MUST be made accessible over the network with NFS and/or SAMBA for OSs to mount.

The reason for this is that applications may live on different servers or virtual machines, but they may still:

- need to use the SQLite database
- need file storage like access
