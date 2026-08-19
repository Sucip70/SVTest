# SVTest API

A small Go REST API for managing posts. It uses MySQL-compatible databases such as TiDB and applies the database migrations automatically when the service starts.

## Requirements

- Go 1.26 or newer
- MySQL or TiDB

## Configuration

Create a `.env` file for local development, or provide the same variables through the deployment platform:

```env
DB_DSN=username:password@tcp(host:3306)/database?charset=utf8mb4&parseTime=True&loc=Local&tls=true
PORT=8080
```

`DB_DSN` is required. `PORT` is optional and defaults to `8080`.

Do not commit `.env` or database credentials to the repository.

## Run Locally

Install dependencies and start the API from the project root:

```bash
go mod download
go run ./cmd/app
```

The API is available at `http://localhost:8080` unless another `PORT` is configured. Migrations in `migrations/` run automatically during startup.

## Endpoints

| Method | Path | Description |
| --- | --- | --- |
| GET | `/health` | Check whether the service is running |
| GET | `/posts` | List posts with filtering and sorting |
| GET | `/post?id=1` | Get one post |
| POST | `/post` | Create a post |
| PUT | `/post?id=1` | Update a post |
| PATCH | `/post/status?id=1` | Update a post status |

### List Posts

`GET /posts` supports these filters:

- `title`: partial match
- `content`: partial match
- `status`: exact match
- `category`: exact match
- `created_date`: date in `YYYY-MM-DD` format

Sorting is controlled by:

- `sort_by`: `title`, `content`, `status`, `category`, or `created_date`
- `sort_order`: `asc` or `desc`
- `page`: page number, starting at `1`
- `limit`: number of posts per page, from `1` to `100`

The default page is `1`, the default limit is `10`, and the default sort is `created_date desc`.

Example:

```text
GET /posts?status=Publish&category=tech&sort_by=title&sort_order=asc
```

The response includes the posts in `data` and the pagination details in `pagination`:

```json
{
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 0,
    "total_pages": 0
  }
}
```

### Create a Post

```http
POST /post
Content-Type: application/json
```

```json
{
  "title": "A post title",
  "content": "Post content",
  "category": "tech",
  "status": "Draft"
}
```

Valid statuses are `Publish`, `Draft`, and `Trash`.

### Update a Post

```http
PUT /post?id=1
Content-Type: application/json
```

```json
{
  "title": "Updated title",
  "content": "Updated content",
  "category": "tech"
}
```

### Update Status

```http
PATCH /post/status?id=1
Content-Type: application/json
```

```json
{
  "status": "Publish"
}
```

## Build

```bash
go build -o app ./cmd/app
```

## Deploy to Cloud Run

You can access the deployment on https://svtest-1014951496037.asia-southeast2.run.app