package id

import (
	"encoding/json"
	"io"
	"time"

	"github.com/frizz925/covid19japan-chatbot/internal/country"
	"github.com/frizz925/covid19japan-chatbot/internal/data"
)

type UpdateResponse struct {
	Update Update `json:"update"`
}

type Update struct {
	Penambahan Penambahan `json:"penambahan"`
	Total      Total      `json:"total"`
}

type Penambahan struct {
	JumlahPositif   int    `json:"jumlah_positif"`
	JumlahMeninggal int    `json:"jumlah_meninggal"`
	JumlahSembuh    int    `json:"jumlah_sembuh"`
	JumlahDirawat   int    `json:"jumlah_dirawat"`
	Tanggal         string `json:"tanggal"`
	Created         string `json:"created"`
}

type Total struct {
	JumlahPositif   int `json:"jumlah_positif"`
	JumlahMeninggal int `json:"jumlah_meninggal"`
	JumlahSembuh    int `json:"jumlah_sembuh"`
	JumlahDirawat   int `json:"jumlah_dirawat"`
}

func ParseUpdate(r io.Reader) (*UpdateResponse, error) {
	var ur UpdateResponse
	err := ur.Parse(r)
	if err != nil {
		return nil, err
	}
	return &ur, nil
}

func (ur *UpdateResponse) Parse(r io.Reader) error {
	return json.NewDecoder(r).Decode(ur)
}

func (ur *UpdateResponse) Normalize() (*data.DailySummary, error) {
	dt, err := parseDateTime(ur.Update.Penambahan.Created)
	if err != nil {
		return nil, err
	}
	return &data.DailySummary{
		Country:             country.INDONESIA,
		CountryID:           country.ID_INDONESIA,
		DateTime:            dt,
		Confirmed:           ur.Update.Penambahan.JumlahPositif,
		Recovered:           ur.Update.Penambahan.JumlahSembuh,
		Deceased:            ur.Update.Penambahan.JumlahMeninggal,
		ConfirmedCumulative: ur.Update.Total.JumlahPositif,
		RecoveredCumulative: ur.Update.Total.JumlahSembuh,
		DeceasedCumulative:  ur.Update.Total.JumlahMeninggal,
	}, nil
}

func parseDateTime(text string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", text)
}