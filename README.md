# Gator CLI

Gator is a command-line RSS feed aggregator written in Go. It allows you to follow RSS feeds, aggregate posts, and browse them directly from your terminal.

## Prerequisites

Before running Gator, ensure you have the following installed:

- **Go**: You need Go installed to compile and install the program. You can download it from [go.dev](https://go.dev/).
- **PostgreSQL**: Gator uses a Postgres database to store user and feed information.

## Installation

To install the Gator CLI, run the following command in your terminal:

```bash
go install github.com/harshaditya-sharma/gator@latest
```

This will compile the project and place the `gator` binary in your `$GOPATH/bin` directory. Make sure this directory is in your system's `PATH`.

## Configuration

Gator requires a configuration file located at `~/.gatorconfig.json`. This file needs to contain the database connection URL and the username of the currently logged-in user.

Create a file named `.gatorconfig.json` in your home directory with the following structure:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": "your_username"
}
```

- `db_url`: The connection string for your Postgres database.
- `current_user_name`: The username you want to operate as (this is updated automatically when you login).

## Usage

Here are some common commands you can use with Gator:

- **Register a new user:**

  ```bash
  gator register <username>
  ```

- **Login as a user:**

  ```bash
  gator login <username>
  ```

- **Add a new RSS feed:**

  ```bash
  gator addfeed <feed_name> <feed_url>
  ```

- **Follow a feed:**

  ```bash
  gator follow <feed_url>
  ```

- **List feeds you are following:**
  ```bash
  gator following
  ```
- **Unfollow a feed:**

  ```bash
  gator unfollow <feed_url>
  ```

- **List all feeds:**

  ```bash
  gator feeds
  ```

- **Aggregate posts (this runs in the background or purely to fetch updates):**

  ```bash
  gator agg <time_between_requests> [concurrency]
  ```

  Arguments:
  - `time_between_requests`: Duration between scrapes (e.g., `1s`, `1m`, `1h`)
  - `concurrency`: Number of concurrent workers (default: 4)

  Example: `gator agg 1m 10`

- **Browse posts:**

  ```bash
  gator browse <limit> <page> [flags]
  ```

  Arguments:
  - `limit`: Number of posts to retrieve (default: 2)
  - `page`: Page number (default: 1)

  Flags:
  - `--sort` or `-s`: Sort order (`asc` or `desc`)
  - `--feed` or `-f`: Filter by feed name

  Example: `gator browse 5 2 --sort asc --feed "Hacker News"`

- **Search posts:**

  ```bash
  gator search <query>
  ```

  Performs a fuzzy search on post titles and descriptions.

- **Like a post:**

  ```bash
  gator like <post_url>
  ```

- **Unlike a post:**

  ```bash
  gator unlike <post_url>
  ```

- **List liked posts:**

  ```bash
  gator liked <limit> <page>
  ```

- **List all users:**

  ```bash
  gator users
  ```

- **Reset the database (clears all users and feeds):**
  ```bash
  gator reset
  ```

## Development

To run the program locally without installing:

```bash
go run . <command> <args>
```
