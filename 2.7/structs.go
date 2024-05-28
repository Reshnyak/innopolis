package main

import "fmt"

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

type Object struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Result struct {
	ObjectID  int `json:"object_id"`
	StudentID int `json:"student_id"`
	Result    int `json:"result"`
}
type ContolSection struct {
	Students []Student
	Objects  []Object
	Results  []Result
}

// Метод возвращающий имя студента по его id
func (cs ContolSection) GetStudentNameByID(id int) string {
	for _, s := range cs.Students {
		if s.ID == id {
			return s.Name
		}
	}
	return ""
}

// Метод возвращающий grade студента по его id
func (cs ContolSection) GetStudentGradeByID(id int) int {
	for _, s := range cs.Students {
		if s.ID == id {
			return s.Grade
		}
	}
	return 0
}

// Метод возвращающий наименование предмета по id
func (cs ContolSection) GetObjectNameByID(id int) string {
	for _, obj := range cs.Objects {
		if obj.ID == id {
			return obj.Name
		}
	}
	return ""
}

//Форматированный вывод =) ...надо попробовать повторить c template
func (cs ContolSection) PrintControlSection() {
	fmt.Println("___________________________________________________")
	fmt.Printf("Student name \t| Grade\t| Object\t| Resulte |\n")
	fmt.Println("___________________________________________________")
	for _, r := range cs.Results {
		fmt.Printf("%s \t\t| ", cs.GetStudentNameByID(r.StudentID))
		fmt.Printf(" %d\t| ", cs.GetStudentGradeByID(r.StudentID))
		objName := cs.GetObjectNameByID(r.ObjectID)
		if len(objName) < 5 {
			fmt.Printf(" %s\t\t| ", cs.GetObjectNameByID(r.ObjectID))
		} else {
			fmt.Printf(" %s\t| ", cs.GetObjectNameByID(r.ObjectID))
		}
		fmt.Printf(" %d\t  |\n", r.Result)
	}

}
func (cs ContolSection) FindTopStudent() []Result {
	tops := Filter(cs.Results, func(result Result) bool {
		stud := Filter(cs.Results, func(r Result) bool {
			if r.StudentID == result.StudentID {
				return true
			}
			return false
		})
		if sum := Reduce(stud, 0, func(a Result, b int) int {
			return a.Result + b
		}); sum == len(cs.Objects)*5 {
			return true
		}
		return false
	})
	return tops
}
func (cs ContolSection) PrintTopStudent() {
	fmt.Println("___________________________________________________")
	fmt.Printf("Student name \t| Grade\t| Object\t| Resulte |\n")
	fmt.Println("___________________________________________________")
	res := cs.FindTopStudent()
	for _, r := range res {
		fmt.Printf("%s \t\t| ", cs.GetStudentNameByID(r.StudentID))
		fmt.Printf(" %d\t| ", cs.GetStudentGradeByID(r.StudentID))
		objName := cs.GetObjectNameByID(r.ObjectID)
		if len(objName) < 5 {
			fmt.Printf(" %s\t\t| ", cs.GetObjectNameByID(r.ObjectID))
		} else {
			fmt.Printf(" %s\t| ", cs.GetObjectNameByID(r.ObjectID))
		}
		fmt.Printf(" %d\t  |\n", r.Result)
	}

}
