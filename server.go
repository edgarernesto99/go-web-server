package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var subjects map[string]map[string]float64 = make(map[string]map[string]float64)
var students map[string]map[string]float64 = make(map[string]map[string]float64)

type Args struct {
	Subject string
	Student string
	Grade   float64
}

//Funciones http
func index(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("index.html"),
	)
}

func formSetGrade(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("formSetGrade.html"),
	)
}

func setGrade(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		//fmt.Println(req.PostForm)
		grade, _ := strconv.ParseFloat(req.FormValue("grade"), 64)
		args := Args{Subject: req.FormValue("subject"), Student: req.FormValue("student"), Grade: grade}
		err := saveGrade(args)
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		if err != nil {
			fmt.Fprintf(
				res,
				cargarHtml("error.html"),
				err,
			)
		} else {
			fmt.Fprintf(
				res,
				cargarHtml("gradeSaved.html"),
				"La calificación de "+args.Student+" en la materia de "+args.Subject+" se registro con exito",
			)
		}
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("tabla.html"),
			studentsToString(),
			subjectsToString(),
		)
	}
}

//Funciones auxiliares
func saveGrade(args Args) error {
	_, subjectExists := subjects[args.Subject]
	_, isStudent := subjects[args.Subject][args.Student]
	_, studentExists := students[args.Student]
	if subjectExists && isStudent {
		return errors.New("El estudiante ya tiene calificación en " + args.Subject)
	}
	if !subjectExists {
		newStudent := make(map[string]float64)
		newStudent[args.Student] = args.Grade
		subjects[args.Subject] = newStudent
	} else {
		subjects[args.Subject][args.Student] = args.Grade
	}
	if !studentExists {
		newSubject := make(map[string]float64)
		newSubject[args.Subject] = args.Grade
		students[args.Student] = newSubject
	} else {
		students[args.Student][args.Subject] = args.Grade
	}
	return nil
}

//Funciones complementarias
func cargarHtml(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func studentsToString() string {
	var html string
	for student := range students {
		html += "<h2>" + student + "</h2>"
		html += "<table border='1'><tr> <th>Materia</th> <th>Calificación</th> </tr>"
		for subject, grade := range students[student] {
			html += "<tr>" +
				"<td>" + subject + "</td>" +
				"<td>" + fmt.Sprint(grade) + "</td>" +
				"</tr>"
		}
		html += "</table>"
	}
	return html
}

func subjectsToString() string {
	var html string
	for subject := range subjects {
		html += "<h2>" + subject + "</h2>"
		html += "<table border='1'><tr> <th>Estudiante</th> <th>Calificación</th> </tr>"
		for student, grade := range subjects[subject] {
			html += "<tr>" +
				"<td>" + student + "</td>" +
				"<td>" + fmt.Sprint(grade) + "</td>" +
				"</tr>"
		}
		html += "</table>"
	}
	return html
}

func main() {
	http.HandleFunc("/root", index) //Pagina de inicio
	http.HandleFunc("/form-set-grade", formSetGrade) //Formulario agregar calificacion
	http.HandleFunc("/set-grade", setGrade) //Guarda la calificacion y muestra si se completo
	fmt.Println("Corriendo servirdor...")
	http.ListenAndServe(":9000", nil)
}