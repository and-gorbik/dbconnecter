# dbconnecter

Простой подключатор к нескольким репликам бд.

Пример использования:
```go
replica1, err := sql.Open(...)
replica2, err := sql.Open(...)
replica3, err := sql.Open(...)
...
connecter := NewMultipleDBConnecter(replica1, replica2, replica3)

conn, err := connecter.Connection()
if err != nil {
    log.Fatal(err)
}

_, err := conn.Exec("...")
...
```