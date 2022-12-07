package model

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
)

type Appointment struct {
	AppointmentID string `json:"AppointmentID"`
	StartTime     string `json:"StarTime"`
	PatientID     string `json:"PatientID"`
	PatientName   string `json:"PatientName"`
}

type Appointments struct {
	Bookings  []*Appointment
	ErrorLog  *log.Logger
	Timeslots []string
}

func (a *Appointments) LoadAppointments() error {
	file, err := os.ReadFile("internal/model/appointments.json")
	if err != nil {
		a.ErrorLog.Print("Failed to open file")
		return err
	}

	var bookings []*Appointment
	if err = json.Unmarshal(file, &bookings); err != nil {
		a.ErrorLog.Print("Failed to unmarshal appointments")
		return err
	}
	a.Timeslots = []string{"00:09:00", "00:10:00", "00:11:00", "00:12:00", "00:13:00", "00:14:00", "00:15:00", "00:16:00", "00:17:00"}
	a.Bookings = bookings
	a.LoadTimeslots()
	return nil
}

func (a *Appointments) LoadTimeslots() {
	for _, appointment := range a.Bookings {
		for i, time := range a.Timeslots {
			if appointment.StartTime == time {
				a.Timeslots = append(a.Timeslots[:i], a.Timeslots[i+1:]...)
				break
			}
		}
	}

}

func (a *Appointments) CreateAppointment(patientID, patientName, startTime string) (string, error) {
	a.LoadAppointments()
	appointmentID := uuid.NewString()

	newAppointment := &Appointment{
		AppointmentID: appointmentID,
		StartTime:     startTime,
		PatientID:     patientID,
		PatientName:   patientName,
	}

	a.Bookings = append(a.Bookings, newAppointment)

	// Write to appointments json
	file, _ := json.MarshalIndent(a.Bookings, "", " ")
	err := os.WriteFile("internal/model/appointments.json", file, 0644)
	if err != nil {
		a.ErrorLog.Print("Writing to file failed. Appointment craetion")
		return "", err
	}
	return appointmentID, nil
}

func (a *Appointments) GetAppointment(appointmentID string) (*Appointment, error) {
	for _, appointment := range a.Bookings {
		if appointment.AppointmentID == appointmentID {
			return appointment, nil
		}
	}
	err := errors.New("Appointment not found")

	return nil, err
}

func (a *Appointments) UserAppointments(patientID string) []*Appointment {
	a.LoadAppointments()
	a.LoadTimeslots()
	var result []*Appointment
	for _, appointment := range a.Bookings {
		if appointment.PatientID == patientID {
			result = append(result, appointment)
		}
	}

	return result
}

func (a *Appointments) adminAppointments() []*Appointment {
	var result []*Appointment
	for _, appointment := range a.Bookings {
		result = append(result, appointment)
	}

	return result
}

func (a *Appointments) UpdateAppointment(appointmentID, time string) error {
	var toUpdate *Appointment
	a.LoadAppointments()

	for _, appointment := range a.Bookings {
		if appointment.AppointmentID == appointmentID {
			toUpdate = appointment
			break
		}
	}
	toUpdate.StartTime = time

	// Write to file update
	file, _ := json.MarshalIndent(a.Bookings, "", " ")
	err := os.WriteFile("internal/model/appointments.json", file, 0644)
	if err != nil {
		a.ErrorLog.Print("Writing to file failed. Appointment update")
		return err
	}
	err = a.ReplaceTimeslot(time)
	if err != nil {
		return err
	}
	return nil

}

func (a *Appointments) DeleteAppointment(appointmentID string) error {
	for i, appointment := range a.Bookings {
		if appointment.AppointmentID == appointmentID {
			a.Bookings = append(a.Bookings[:i], a.Bookings[i+1:]...)
			break
		}
	}

	err := errors.New("Failed to delete. Apopintment not found")

	// Write to appointments json
	file, _ := json.MarshalIndent(a.Bookings, "", " ")
	err = os.WriteFile("internal/model/appointments.json", file, 0644)
	if err != nil {
		a.ErrorLog.Print("Writing to file failed. Appointment craetion")
		return err
	}
	return nil
}

func (a *Appointments) ReplaceTimeslot(time string) error {
	index, err := strconv.Atoi(time[3:5])
	if err != nil {
		a.ErrorLog.Print("Failed to convert")
		return err
	}

	index -= 9
	a.Timeslots = append(a.Timeslots[:index+1], a.Timeslots[index:]...)
	a.Timeslots[index] = time
	return nil

}
