package main

//Тут не много сложнее. Не совсем понял задания поэтому сделаю два варианта

// Вариант первый оборачивает текст в несколько форматов: _`**текст**`_
// т.е. к тексту по цепочке применяются разные форматы
type ChainFormater_1 struct {
	formats []Formatter
	text    string
}

//Метод добавления формата
func (chf *ChainFormater_1) AddFormatter(f Formatter) {
	chf.formats = append(chf.formats, f)
}

//Форматирование
func (chf *ChainFormater_1) Format() string {
	res := chf.text
	for _, f := range chf.formats {
		switch f.(type) {
		case BoltText:
			res = BoltText{text: res}.Format()
		case CodeText:
			res = CodeText{text: res}.Format()
		case ItalicText:
			res = ItalicText{text: res}.Format()
		default:
			res = PlainText{text: res}.Format()

		}
	}
	return res
}

// Вариант второй выводит несколько форматов текстов: текст**текст1**_текст2_`текст3`
// т.е. цепочка форматов
type ChainFormater_2 struct {
	formats []Formatter
}

//Метод добавления формата
func (chf *ChainFormater_2) AddFormatter(f Formatter) {
	chf.formats = append(chf.formats, f)
}

//Форматирование
func (chf *ChainFormater_2) Format() string {
	res := ""
	for _, f := range chf.formats {
		res += f.Format()
	}
	return res
}
