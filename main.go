package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"io"
	"strconv"
	"sort"
)

const (
	sourceLink	= "http://82.196.1.83:9570"
	shortForm	= "02.01.2006 15:04"
//const data = "[{\"Last name\": \"John\",\"First name\": \"Doe\",\"Age\": \"35\",\"Gender\": \"m\",\"marital\": \"true\",\"Last login\": \"14.05.2016 12:37\"},{\"Last name\": \"Joshua\",\"First name\": \"Star\",\"Age\": \"17\",\"Gender\": \"m\",\"marital\": \"false\",\"Last login\": \"12.05.2016 17:30\"},{\"Last name\": \"Jane\",\"First name\": \"Sith\",\"Age\": \"20\",\"Gender\": \"f\",\"marital\": \"true\",\"Last login\": \"02.05.2016 07:22\"},{\"Last name\": \"Robert\",\"First name\": \"Milson\",\"Age\": \"44\",\"Gender\": \"m\",\"marital\": \"true\",\"Last login\": \"10.04.2016 10:00\"},{\"Last name\": \"Elisabeth\",\"First name\": \"Morth\",\"Age\": \"19\",\"Gender\": \"f\",\"marital\": \"true\",\"Last login\": \"14.04.2016 17:31\"},{\"Last name\": \"Bary\",\"First name\": \"Sorm\",\"Age\": \"20\",\"Gender\": \"m\",\"marital\": \"false\",\"Last login\": \"11.01.2016 12:27\"},{\"Last name\": \"Mary\",\"First name\": \"Douson\",\"Age\": \"33\",\"Gender\": \"f\",\"marital\": \"true\",\"Last login\": \"\"},{\"Last name\": \"Sarah\",\"First name\": \"Connor\",\"Age\": \"28\",\"Gender\": \"f\",\"marital\": \"true\",\"Last login\": \"02.03.2014 22:31\"},{\"Last name\": \"Alistar\",\"First name\": \"Tampler\",\"Age\": \"30\",\"Gender\": \"m\",\"marital\": \"true\",\"Last login\": \"12.01.2016 17:30\"},{\"Last name\": \"Veronica\",\"First name\": \"Ingil\",\"Age\": \"22\",\"Gender\": \"f\",\"marital\": \"false\",\"Last login\": \"22.04.2016 12:55\"},{\"Last name\": \"Martin\",\"First name\": \"Mulen\",\"Age\": \"12\",\"Gender\": \"m\",\"marital\": \"false\",\"Last login\": \"18.06.2015 14:00\"},{\"Last name\": \"Monica\",\"First name\": \"Doust\",\"Age\": \"76\",\"Gender\": \"f\",\"marital\": \"true\",\"Last login\": \"\"},{\"Last name\": \"Michael\",\"First name\": \"Kurst\",\"Age\": \"27\",\"Gender\": \"m\",\"marital\": \"false\",\"Last login\": \"07.03.2015 12:12\"}]"
)

type InputData struct {
	Last_name	string
	First_name	string
	Age		int
	Gender		string
	Marital		bool
	Last_login	time.Time
}

var (
	_ json.Unmarshaler = (*InputData)(nil)
)

func (f *InputData)UnmarshalJSON(data []byte) error {
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		fmt.Println(err)
		return err
	}
	f.Last_name	= m["Last name"]
	f.First_name	= m["First name"]
	f.Age, _	= strconv.Atoi(m["Age"])
	f.Gender 	= m["Gender"]
	f.Marital, _	= strconv.ParseBool(m["marital"])
	f.Last_login, _	= time.Parse(shortForm, m["Last login"])
	return nil
}

func main() {

	resp, err := getHTTPData(sourceLink)
	if err != nil {
		return
	}

//	resp := strings.NewReader(data)
	userdata, err := getRawData(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	SortByAge(userdata)

	compileData(userdata)
	return
}

func getRawData(r io.Reader) (userdata []InputData, err error) {
	err = json.NewDecoder(r).Decode(&userdata)
	return userdata, err
}

func getHTTPData(l string) (io.Reader, error) {
	resp, err := http.Get(l)
	if err != nil || resp == nil {
		fmt.Println(err)
		return nil, err
	}
	return resp.Body, err
}

func (s InputData)String() string {
	return fmt.Sprintf("Last name = %s, First name = %s, Age = %d, Gender = %s, Marital = %t, Last login = %s, Last login Unix= %d\n",
	s.Last_name, s.First_name, s.Age, s.Gender, s.Marital, s.Last_login.Format(time.UnixDate), s.Last_login.Unix())
}

func SortByAge(s []InputData) {
	ageUp := func(p1, p2 *InputData) bool {
		return p1.Age < p2.Age
	}
	ageDown := func(p1, p2 *InputData) bool {
		return p1.Age > p2.Age
	}
	By(ageUp).Sort(s)
	fmt.Println("Sort by age up:\n", s)

	By(ageDown).Sort(s)
	fmt.Println("Sort by age down:\n", s)
}

type By func(p1, p2 *InputData) bool

func (by By) Sort(data []InputData) {
	ps := &InputDataSorter{
		data: data,
		by:     by,
	}
	sort.Sort(ps)
}

type InputDataSorter struct {
	data []InputData
	by      func(p1, p2 *InputData) bool
}

func (s *InputDataSorter) Len() int {
	return len(s.data)
}

func (s *InputDataSorter) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

func (s *InputDataSorter) Less(i, j int) bool {
	return s.by(&s.data[i], &s.data[j])
}

func compileData(d []InputData) {
	type person struct {
		LastName	string
		FirstName	string
	}
	ar := make(map[string]map[bool][]person)
	ar["m"] = make(map[bool][]person)
	ar["f"] = make(map[bool][]person)

	for _, el := range d {
		ar[el.Gender][el.Marital] = append(ar[el.Gender][el.Marital], person{el.Last_name, el.First_name})
	}

	fmt.Println("Compiled data:\n", ar)

	return
}



