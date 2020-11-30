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

func formGetStudentAverage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("formGetStudentAverage.html"),
		studentsToSelect(),
	)
}

func showStudentAverage(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(res, "ParseForm() error %v", err)
		return
	}
	student := req.FormValue("student")
	var average float64
	err := getStudentAverage(student, &average)
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
			cargarHtml("showAverage.html"),
			"Promedio de "+student+": "+fmt.Sprint(average),
		)
	}
}

func formGetSubjectAverage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("formGetSubjectAverage.html"),
		subjectsToSelect(),
	)
}

func showSubjectAverage(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(res, "ParseForm() error %v", err)
		return
	}
	subject := req.FormValue("subject")
	var average float64
	err := getSubjectAverage(subject, &average)
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
			cargarHtml("showAverage.html"),
			"Promedio en "+subject+": "+fmt.Sprint(average),
		)
	}
}

func showGeneralAverage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	var average float64
	err := getGeneralAverage(&average)
	if err != nil {
		fmt.Fprintf(
			res,
			cargarHtml("error.html"),
			err,
		)
	} else {
		fmt.Fprintf(
			res,
			cargarHtml("showAverage.html"),
			"Promedio general: "+fmt.Sprint(average),
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

func getStudentAverage(student string, reply *float64) error {
	_, studentExists := students[student]
	if !studentExists {
		return errors.New("El alumno no existe")
	}
	var average float64 = 0
	var total float64 = 0
	for _, grade := range students[student] {
		average = average + grade
		total = total + 1
	}
	average = average / total
	*reply = average
	return nil
}

func getSubjectAverage(subject string, reply *float64) error {
	_, subjectExists := subjects[subject]
	if !subjectExists {
		return errors.New("La materia no existe")
	}
	var average float64 = 0
	var total float64 = 0
	for _, grade := range subjects[subject] {
		average = average + grade
		total = total + 1
	}
	average = average / total
	*reply = average
	return nil
}

func getGeneralAverage(reply *float64) error {
	var generalAverage float64 = 0
	var generalTotal float64 = 0
	for student := range students {
		var average float64 = 0
		total := 0
		generalTotal = generalTotal + 1
		for _, grade := range students[student] {
			average = average + grade
			total = total + 1
		}
		average = average / float64(total)
		generalAverage = generalAverage + average
	}
	if generalTotal == 0 {
		return errors.New("No hay estudiantes")
	}
	generalAverage = generalAverage / generalTotal
	*reply = generalAverage
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

func studentsToSelect() string {
	var html string
	for student := range students {
		html += "<option value='" + student + "' id='student'>" + student + "</option>"
	}
	return html
}

func subjectsToSelect() string {
	var html string
	for subject := range subjects {
		html += "<option value='" + subject + "' id='subject'>" + subject + "</option>"
	}
	return html
}

func main() {
	http.HandleFunc("/root", index) //Pagina de inicio
	http.HandleFunc("/form-set-grade", formSetGrade) //Formulario agregar calificacion
	http.HandleFunc("/set-grade", setGrade) //Guarda la calificacion y muestra si se completo
	http.HandleFunc("/form-student-average", formGetStudentAverage) // Formulario obtener prom de estudiante
	http.HandleFunc("/form-subject-average", formGetSubjectAverage) // Formulario obtener prom de materia
	http.HandleFunc("/student-average", showStudentAverage) //Muestra promedio de estudiante
	http.HandleFunc("/subject-average", showSubjectAverage) //Mostrar promedio de materia
	http.HandleFunc("/general-average", showGeneralAverage) //Mostrar promedio general
	fmt.Println("Corriendo servirdor...")
	http.ListenAndServe(":9000", nil)
}