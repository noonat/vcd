package vcd

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeVehicleData(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "Valid data should not be modified",
			in:   `<span style="color: rgb(255,0,0)">&nbsp;<br/>&nbsp;</span>&nbsp;`,
			out:  `<span style="color: rgb(255,0,0)">&nbsp;<br/>&nbsp;</span>&nbsp;`,
		},
		{
			name: "Valid data inside an invalid element should be ignored",
			in:   `<div><span style="color: rgb(255,0,0)">&nbsp;<br/>&nbsp;</span>&nbsp;</div>`,
			out:  ``,
		},
		{
			name: "Valid data outside an invalid element should not be modified",
			in:   `&nbsp;<div><span style="color: rgb(255,0,0)">&nbsp;<br/>&nbsp;</span>&nbsp;</div><br/>`,
			out:  `&nbsp;<br/>`,
		},
		{
			name: "Invalid attributes should be stripped",
			in:   `&nbsp;<span id="foo" style="color: rgb(255,0,0)">&nbsp;<br style="font-weight: bold;"/>&nbsp;</span>&nbsp;`,
			out:  `&nbsp;<span style="color: rgb(255,0,0)">&nbsp;<br/>&nbsp;</span>&nbsp;`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, err := sanitizeVesselData(test.in)
			if assert.NoError(t, err) {
				assert.Equal(t, test.out, out)
			}
		})
	}
}

func TestNewVesselFromString(t *testing.T) {
	s := `<tt style="background-color: rgb(0,0,0)">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span style="color: rgb(0,255,0)">|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">/---\</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;=|&nbsp;&nbsp;=</span><span style="color: rgb(255,51,51)">|&nbsp;@&nbsp;|</span><span style="color: rgb(0,255,0)">=&nbsp;&nbsp;|=&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">|---|</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span></tt><br/><a href="http://www.captainforever.com/captainforever.php?cfe=nqlbqna2lfhwbjuo0w0nbn1cnmbmntagrhxa5undxu">Pilot this vessel</a>`
	v, err := NewVesselFromString(s)
	if assert.NoError(t, err) && assert.NotNil(t, v) {
		assert.Equal(t, *v, Vessel{
			CFE:  "nqlbqna2lfhwbjuo0w0nbn1cnmbmntagrhxa5undxu",
			Data: `&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span style="color: rgb(0,255,0)">|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">/---\</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;=|&nbsp;&nbsp;=</span><span style="color: rgb(255,51,51)">|&nbsp;@&nbsp;|</span><span style="color: rgb(0,255,0)">=&nbsp;&nbsp;|=&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">|---|</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>`,
		})
	}
}

func TestQueryVesselByID(t *testing.T) {
	db := requireTestDB(t, []string{"vessels.sql"})
	defer db.Close()

	tests := []struct {
		id       int
		expected *Vessel
	}{
		{
			id: 1,
			expected: &Vessel{
				ID:        1,
				CreatedAt: time.Date(2009, 9, 9, 4, 11, 0, 0, time.UTC),
				IP:        net.ParseIP("76.115.235.156"),
				CFE:       "nqlbqnbthgnxbyun1w0nnt1cqgbana2lrxwa5uqcxu",
				Data:      `&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span style="color: rgb(0,255,0)">|&nbsp;|---|&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;|&nbsp;|&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">/---\</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;=</span><span style="color: rgb(255,51,51)">|&nbsp;@&nbsp;|</span><span style="color: rgb(0,255,0)">=&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">|---|</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;|&nbsp;|&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|---|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>`,
			},
		},
		{
			id: 2,
			expected: &Vessel{
				ID:        2,
				CreatedAt: time.Date(2009, 9, 9, 4, 11, 2, 0, time.UTC),
				IP:        net.ParseIP("76.115.234.44"),
				CFE:       "nqlbqnnaghnxnoun1wwxnbcedhnettolgmwa2uaswwnbasdmgn5ta3lghwbjuocxzxaqp4qgne",
				Data:      `&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span style="color: rgb(0,255,0)">|&nbsp;&nbsp;&nbsp;</span><span style="color: rgb(255,255,0)">|&nbsp;&nbsp;&nbsp;</span><span style="color: rgb(0,255,0)">|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;v&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;</span><span style="color: rgb(255,255,0)">+&nbsp;&nbsp;&nbsp;</span><span style="color: rgb(0,255,0)">+&nbsp;&nbsp;&nbsp;v&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">/---\</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;=|&nbsp;&nbsp;=</span><span style="color: rgb(255,51,51)">|&nbsp;@&nbsp;|</span><span style="color: rgb(0,255,0)">=&nbsp;&nbsp;|=&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">|---|</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;|&nbsp;|&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|---|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>`,
			},
		},
		{
			id:       1337,
			expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("vessel id %d", test.id), func(t *testing.T) {
			v, err := QueryVesselByID(context.Background(), db, test.id)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expected, v)
			}
		})
	}
}

func TestQueryVessels(t *testing.T) {
	db := requireTestDB(t, []string{"vessels.sql"})
	defer db.Close()

	vessels, err := QueryVessels(context.Background(), db)
	if assert.NoError(t, err) {
		assert.Equal(t, vessels, []*Vessel{
			&Vessel{
				ID:        1,
				CreatedAt: time.Date(2009, 9, 9, 4, 11, 0, 0, time.UTC),
				IP:        net.ParseIP("76.115.235.156"),
				CFE:       "nqlbqnbthgnxbyun1w0nnt1cqgbana2lrxwa5uqcxu",
				Data:      `&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span style="color: rgb(0,255,0)">|&nbsp;|---|&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;|&nbsp;|&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">/---\</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;=</span><span style="color: rgb(255,51,51)">|&nbsp;@&nbsp;|</span><span style="color: rgb(0,255,0)">=&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">|---|</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;|&nbsp;|&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|---|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>`,
			},
			&Vessel{
				ID:        2,
				CreatedAt: time.Date(2009, 9, 9, 4, 11, 2, 0, time.UTC),
				IP:        net.ParseIP("76.115.234.44"),
				CFE:       "nqlbqnnaghnxnoun1wwxnbcedhnettolgmwa2uaswwnbasdmgn5ta3lghwbjuocxzxaqp4qgne",
				Data:      `&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<span style="color: rgb(0,255,0)">|&nbsp;&nbsp;&nbsp;</span><span style="color: rgb(255,255,0)">|&nbsp;&nbsp;&nbsp;</span><span style="color: rgb(0,255,0)">|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;v&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;</span><span style="color: rgb(255,255,0)">+&nbsp;&nbsp;&nbsp;</span><span style="color: rgb(0,255,0)">+&nbsp;&nbsp;&nbsp;v&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">/---\</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;=|&nbsp;&nbsp;=</span><span style="color: rgb(255,51,51)">|&nbsp;@&nbsp;|</span><span style="color: rgb(0,255,0)">=&nbsp;&nbsp;|=&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;---&nbsp;---&nbsp;</span><span style="color: rgb(255,51,51)">|---|</span><span style="color: rgb(0,255,0)">&nbsp;---&nbsp;---&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;|&nbsp;|&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;+&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|---|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;^&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<br/>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>`,
			},
		})
	}
}

func TestQueryVesselsEmpty(t *testing.T) {
	db := requireTestDB(t, nil)
	defer db.Close()

	vessels, err := QueryVessels(context.Background(), db)
	if assert.NoError(t, err) {
		assert.Equal(t, vessels, []*Vessel{})
	}
}
