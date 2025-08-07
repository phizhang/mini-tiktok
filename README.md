# mini tiktok

A mini tiktok backend implementation using golang for microservice and RESTful API,  Cassandra and AWS S3 integration.
I plan to implement te following features:

1. Video upload and list for current user. (done)
2. Live stream server. (Ongoing)
3. Following and unfollowing. (Planned)
4. Video feeds from friends. (Planned)
5. Global video feeds based on recommendation engine. (Planned)

## Folder Structure

- `cmd/` - Main applications for this project
- `internal/api/` - API routing and middleware
- `internal/db/` - Cassandra database access logic
- `internal/s3/` - AWS S3 access logic
- `internal/models/` - Data models
- `internal/handlers/` - HTTP handlers
- `internal/utils/` - Utility functions
- `configs/` - Configuration files
- `scripts/` - Helper scripts (e.g., DB migrations)
- `pkg/` - Public Go packages (if any)

## Getting Started

1. Install dependencies
2. Configure Cassandra and AWS S3 credentials
3. Run the application from `cmd/`

## Storage Settings

By default the uploaded vidro will be in the local disk of the mini-tiktok server, using user UUID as the bucket.
We can also use environment variable to indicate S3 storage will be used.

## Running Cassandra with Docker

You can quickly start a Cassandra instance using Docker with default credentials:

```sh
docker run --name cassandra -p 9042:9042 -e CASSANDRA_PASSWORD_SEEDER=yes -e CASSANDRA_USER=cassandra -e CASSANDRA_PASSWORD=cassandra -d cassandra:latest
```

- Default username: `cassandra`
- Default password: `cassandra`

Wait a few minutes for Cassandra to initialize before connecting.

## Connecting to Cassandra and Creating the `video_meta` Table

1. Connect to Cassandra using cqlsh:

```sh
docker exec -it cassandra cqlsh -u cassandra -p cassandra
```

2. Create a keyspace (if not already created):

```sql
CREATE KEYSPACE IF NOT EXISTS mini_tiktok WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
```

3. Use the keyspace:

```sql
USE mini_tiktok;
```

4. Create the `video_meta` table:

```sql
CREATE TABLE IF NOT EXISTS video_meta (
    user_id UUID,
    video_id UUID,
    title TEXT,
    tag TEXT,
    created_time TIMESTAMP,
    location TEXT,
    PRIMARY KEY ((user_id), video_id, created_time)
);
```

## Example: Upload a Video File via cURL

You can upload a video file using the following cURL command:

```sh
curl -X POST \
  -H "X-User-ID: <your_user_id>" \
  -F "title=My Video Title" \
  -F "tag=funny" \
  -F "file=@/path/to/your/video.mp4" \
  http://localhost:8080/api/upload
```

Replace `<your_user_id>` with your user ID and `/path/to/your/video.mp4` with the path to your video file.
