package main

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
type ControlSection struct {
	Students []Student
	Objects  []Object
	Results  []Result
}

// Метод возвращающий мапу студентов с ID - key
func (cs ControlSection) GetStudentsMap() map[int]Student {
	res := make(map[int]Student)
	for _, stud := range cs.Students {
		res[stud.ID] = stud
	}
	return res
}

// Метод возвращающий мапу предметов с ID - key
func (cs ControlSection) GetObjectsMap() map[int]Object {
	res := make(map[int]Object)
	for _, obj := range cs.Objects {
		res[obj.ID] = obj
	}
	return res
}
