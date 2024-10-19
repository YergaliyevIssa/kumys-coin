package doctors

import (
	"fmt"
	"os"
)

type Doctor struct {
	Name         string
	DoctorRole   string
	Experience   string
	ServicePlace string
	Price        string
	Clinic       string
	Address      string
	Phone        string
	PhotoURL     string
}

func (d Doctor) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		d.Name,
		d.DoctorRole,
		d.Experience,
		d.ServicePlace,
		d.Price,
		d.Clinic,
		d.Address,
		d.Phone,
	)
}

var Doctors = []Doctor{
	{
		Name:         "Сейтенов Тимур Женисович",
		DoctorRole:   "Терапевт",
		Experience:   "Стаж 12 лет / Врач первой категории",
		ServicePlace: "Прием в клинике",
		Price:        "10000 тг.",
		Clinic:       "Ayala",
		Address:      "ул. Туркестан, 28/2, Есильский район, Левый берег, Астана",
		Phone:        "+7 708 515 7812",
		PhotoURL:     "file://" + os.Getenv("PWD") + "pkg/doctors/dr_timur.jpg",
	},
	{
		Name:         "Елтай Айгерім",
		DoctorRole:   "Терапевт",
		Experience:   "Стаж 9 лет / Врач высшей категории",
		ServicePlace: "Прием в клинике",
		Price:        "8000 тг.",
		Clinic:       "INFINITY LIFE",
		Address:      "проспект Кабанбай батыр, 40, Есильский район, Левый берег, Астана",
		Phone:        "+7 777 651 45 55",
		PhotoURL:     "file://" + os.Getenv("PWD") + "pkg/doctors/dr_aigerim.jpg",
	},
}
