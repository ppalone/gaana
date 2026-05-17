# Gaana

Unofficial Go client for Gaana

## Installation

```bash
go get github.com/ppalone/gaana
```

## Usage

Search Songs

```go
client := gaana.NewClient(nil)
results, err := client.SearchSongs(context.Background(), "animals")
if err != nil {
    log.Fatal(err)
}
fmt.Println(results)
```

Get Song Details by Track ID

```go
client := gaana.NewClient(nil)
song, err := client.GetSongDetailByTrackId(context.Background(), 1783362)
if err != nil {
    log.Fatal(err)
}
fmt.Println(song)
```

Get Song Details by SEO Key

```go
client := gaana.NewClient(nil)
song, err := client.GetSongDetailBySeoKey(context.Background(), "animals-11")
if err != nil {
    log.Fatal(err)
}
fmt.Println(song)
```

## Author 

Pranjal 

## LICENSE

MIT
