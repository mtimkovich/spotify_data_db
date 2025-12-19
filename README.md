# Spotify Listening Data

Convert Spotify extended streaming history into an SQLite DB.

For me, my Spotify listening history was 188 MB of JSON files. This script converts that data into a 22 MB SQLite DB that's much easier and faster to analyze.

Some sample SQL queries are included in the `sql/` directory.

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
