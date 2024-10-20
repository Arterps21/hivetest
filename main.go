package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    _ "github.com/lib/pq"
)

const (
    host     = "postgres"
    port     = 5432
    user     = "postgres"
    password = "214365"
    dbname   = "hive"
)

func main() {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            r.ParseForm()
            message := r.FormValue("message")

            // Вставляем сообщение в базу данных
            _, err = db.Exec("INSERT INTO greetings (message) VALUES ($1)", message)
            if err != nil {
                http.Error(w, "Database error", http.StatusInternalServerError)
                return
            }

            // Отправка сообщения о добавлении и форма для "More"
            fmt.Fprint(w, `
                <html>
                    <body>
                        <h1>Сообщение добавлено</h1>
                        <p>Message added!</p>
                        <form method="get" action="/">
                            <button type="submit">More</button>
                        </form>
                    </body>
                </html>
            `, message)
            return
        }

        // Отправка HTML-формы
        renderForm(w, db)
    })

    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderForm(w http.ResponseWriter, db *sql.DB) {
    fmt.Fprint(w, `
        <html>
            <body>
                <h1>Добавить сообщение</h1>
                <form method="post" action="/">
                    <label for="message">Сообщение:</label>
                    <input type="text" id="message" name="message" required>
                    <button type="submit">Отправить</button>
                </form>
                <h2>Последние сообщения</h2>
                <ul>
    `)

    // Получаем и отображаем последние сообщения
    rows, err := db.Query("SELECT message FROM greetings ORDER BY id DESC LIMIT 5")
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var greeting string
        if err := rows.Scan(&greeting); err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "<li>%s</li>", greeting)
    }

    fmt.Fprint(w, `
                </ul>
            </body>
        </html>
    `)
}
