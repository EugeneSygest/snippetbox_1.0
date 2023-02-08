package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showPost = Article{}

//var database *sql.DB

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./site/ui/html/index.html", "./site/ui/html/header.html", "./site/ui/html/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	// Выборка данных
	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next() {
		var post Article
		err := res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
		// fmt.Println(fmt.Sprintf("Post: %s with id: %d", post.Title, post.Id))
	}
	defer res.Close()

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./site/ui/html/create.html", "./site/ui/html/header.html", "./site/ui/html/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены!")
	} else {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))

		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func Show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //получение id при помощи gorillamux

	//w.WriteHeader(http.StatusOK)

	t, err := template.ParseFiles("./site/ui/html/show.html", "./site/ui/html/header.html", "./site/ui/html/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Выборка данных

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for res.Next() {
		var post Article
		err := res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		showPost = post

	}
	defer res.Close()

	t.ExecuteTemplate(w, "show", showPost)
}

func delete_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //получение id при помощи gorillamux

	t, err := template.ParseFiles("./site/ui/html/delete.html", "./site/ui/html/header.html", "./site/ui/html/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.Exec(fmt.Sprintf("DELETE FROM `articles` WHERE `id` = '%s'", vars["id"]))

	t.ExecuteTemplate(w, "delete", nil)

}

// возвращаем пользователю страницу для редактирования объекта
func edit_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //получение id при помощи gorillamux

	//w.WriteHeader(http.StatusOK)

	t, err := template.ParseFiles("./site/ui/html/edit.html", "./site/ui/html/header.html", "./site/ui/html/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Выборка данных
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	showPost = Article{}
	for res.Next() {
		var post Article
		err := res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		showPost = post

	}
	defer res.Close()

	t.ExecuteTemplate(w, "edit", showPost)
}

func edit_handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")
	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены!")
	} else {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		//insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, full_text))
		insert, err := db.Query(fmt.Sprintf("update `articles` set `title`='%s', `anons`='%s', `full_text` = '%s' where `id` = '%s'", title, anons, full_text, vars["id"]))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func handlefunc() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", Show_post).Methods("GET")
	rtr.HandleFunc("/delete/{id:[0-9]+}", delete_post).Methods("GET")
	rtr.HandleFunc("/edit/{id:[0-9]+}", edit_post).Methods("GET")
	rtr.HandleFunc("/edit/{id:[0-9]+}", edit_handler).Methods("POST")

	http.Handle("/", rtr)

	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.ListenAndServe(":8080", nil)
}

func main() {
	handlefunc()
}
