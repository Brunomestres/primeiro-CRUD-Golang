package servidor

import (
	db "banco-de-dados/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha na requisição!"))
		return
	}

	var usuario usuario

	if erro = json.Unmarshal(body, &usuario); erro != nil {
		w.Write([]byte("Erro ao tranformar o usuario para struct"))
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar o banco"))
		return
	}

	defer db.Close()

	statemant, erro := db.Prepare("insert into usuarios (nome, email) values (?,?)")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	defer statemant.Close()

	insert, erro := statemant.Exec(usuario.Nome, usuario.Email)
	if erro != nil {
		w.Write([]byte("Erro ao criar o usuario"))
		return
	}

	id, erro := insert.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao buscar o id"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("usuario inserido com sucesso Id: %d", id)))
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, erro := db.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar o banco"))
		return
	}
	defer db.Close()

	linhas, erro := db.Query("select * from usuarios")
	if erro != nil {
		w.Write([]byte("Erro ao buscar usuarios"))
		return
	}
	defer linhas.Close()

	var usuarios []usuario
	for linhas.Next() {
		var usuario usuario
		if erro := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear usuarios"))
			return
		}

		usuarios = append(usuarios, usuario)

	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil {
		w.Write([]byte("Erro ao converter usuarios para json"))
		return
	}
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter paramentro"))
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco"))
		return
	}

	linha, erro := db.Query("select * from usuarios where id = ?", ID)
	if erro != nil {
		w.Write([]byte("Erro ao buscar o usuario"))
		return
	}

	var usuario usuario
	if linha.Next() {
		if erro := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear o usuario"))
			return
		}

		w.WriteHeader(http.StatusOK)
		if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
			w.Write([]byte("Erro ao converter usuario para json"))
			return
		}
	}
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter paramentro"))
		return
	}

	body, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha na requisição!"))
		return
	}

	var usuario usuario

	if erro = json.Unmarshal(body, &usuario); erro != nil {
		w.Write([]byte("Erro ao tranformar o usuario para struct"))
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	defer statement.Close()
	if _, erro := statement.Exec(usuario.Nome, usuario.Email, ID); erro != nil {
		w.Write([]byte("Erro ao atualizar o usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter paramentro"))
		return
	}

	db, erro := db.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao deletar o usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
