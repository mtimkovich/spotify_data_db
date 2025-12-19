# Spotify Listening Data

Convert Spotify extended streaming history into an SQLite DB.

Also includes some SQL queries for analyzing the database.

## Usage

To start, you'll need to [request your extended listening history](https://www.spotify.com/us/account/privacy/) from Spotify.

```
go build spotify.go
./spotify [path to listening data]
```

This will create a `spotify.db` file in the current directory.

## Running SQLite queries

```
sqlite3 spotify.db < sql/top_artists_2025.sql
```
