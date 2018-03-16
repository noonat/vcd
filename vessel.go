package vcd

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

var (
	ltRegexp     = regexp.MustCompile(`&lt([^;])`)
	ltReplace    = `&lt;$1`
	vesselRegexp = regexp.MustCompile(`<tt style="background-color: rgb\(0,0,0\)">(.+)</tt><br/><a href="http://www\.captainforever\.com/captainforever\.php\?cfe=([a-z0-9]+)">Pilot this vessel</a>`)
)

// Vessel represents a Captain Forever vessel that someone has saved to VCD.
type Vessel struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IP        net.IP    `json:"ip"`
	CFE       string    `json:"cfe"`
	Data      string    `json:"data"`

	// read-only join fields
	PilotClicks int `json:"pilot_clicks"`
}

// NewVesselFromString creates a new Vessel struct from a Captain Forever
// exported HTML snippet.
func NewVesselFromString(s string) (*Vessel, error) {
	matches := vesselRegexp.FindStringSubmatch(s)
	if matches == nil {
		return nil, errors.New("string didn't match regular expression")
	}

	data, err := sanitizeVesselData(matches[1])
	if err != nil {
		return nil, err
	}

	return &Vessel{CFE: matches[2], Data: data}, nil
}

// sanitizeVesselData parses the HTML from the Captain Forever export and
// only keeps the parts we want. Specifically, we only want to keep span and
// br element nodes and text nodes. On the span nodes, the only attribute that
// is allowed is "style".
func sanitizeVesselData(data string) (string, error) {
	var (
		buf           bytes.Buffer
		skipping      bool
		skippingCount int
	)
	tokenizer := html.NewTokenizer(strings.NewReader(data))
	for tokenizer.Next() != html.ErrorToken {
		raw := string(tokenizer.Raw())
		token := tokenizer.Token()
		switch token.Type {
		case html.StartTagToken:
			if skipping {
				skippingCount++
				continue
			}
			if token.Data != "span" {
				skipping = true
				skippingCount++
				continue
			}
			var attrs []html.Attribute
			for _, attr := range token.Attr {
				if attr.Key != "style" {
					continue
				}
				attrs = append(attrs, attr)
			}
			token.Attr = attrs
			buf.WriteString(token.String())

		case html.EndTagToken:
			if skipping {
				skippingCount--
				if skippingCount == 0 {
					skipping = false
				}
				continue
			}
			buf.WriteString(token.String())

		case html.SelfClosingTagToken:
			if skipping || token.Data != "br" {
				continue
			}
			token.Attr = nil
			buf.WriteString(token.String())

		case html.TextToken:
			if skipping {
				continue
			}
			buf.WriteString(raw)
		}
	}
	if err := tokenizer.Err(); err != io.EOF {
		return "", errors.Wrap(err, "error tokenizing")
	}
	return buf.String(), nil
}

// CreateVessel inserts a new vessel into the database.
func CreateVessel(ctx context.Context, db *sql.DB, ip net.IP, cfe, data string) (*Vessel, error) {
	result, err := db.ExecContext(ctx, `
		INSERT INTO vessels (ip, cfe, data) VALUES (INET_ATON(?), ?, ?)`,
		ip.String(), cfe, data)
	if err != nil {
		return nil, errors.Wrap(err, "error inserting vessel")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "error getting last insert id")
	}
	return QueryVesselByID(ctx, db, int(id))
}

// QueryVesselByCFE finds an existing vessel by CFE.
func QueryVesselByCFE(ctx context.Context, db *sql.DB, cfe string) (*Vessel, error) {
	var (
		ip string
		v  Vessel
	)
	row := db.QueryRowContext(ctx, `
		SELECT v.id, v.created_at, INET_NTOA(v.ip), v.cfe, v.data, COALESCE(vpc.clicks, 0)
		FROM vessels v
		LEFT JOIN (
			SELECT vessel_id, COUNT(*) AS clicks
			FROM vessel_pilot_clicks
			GROUP BY vessel_id
		) vpc ON vpc.vessel_id = v.id
		WHERE v.cfe = ?
	`, cfe)
	if err := row.Scan(&v.ID, &v.CreatedAt, &ip, &v.CFE, &v.Data, &v.PilotClicks); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "error querying for vessel")
	}
	v.IP = net.ParseIP(ip)
	return &v, nil
}

// QueryVesselByID finds an existing vessel by ID.
func QueryVesselByID(ctx context.Context, db *sql.DB, id int) (*Vessel, error) {
	var (
		ip string
		v  Vessel
	)
	row := db.QueryRowContext(ctx, `
		SELECT v.id, v.created_at, INET_NTOA(v.ip), v.cfe, v.data, COALESCE(vpc.clicks, 0)
		FROM vessels v
		LEFT JOIN (
			SELECT vessel_id, COUNT(*) AS clicks
			FROM vessel_pilot_clicks
			GROUP BY vessel_id
		) vpc ON vpc.vessel_id = v.id
		WHERE v.id = ?
	`, id)
	if err := row.Scan(&v.ID, &v.CreatedAt, &ip, &v.CFE, &v.Data, &v.PilotClicks); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "error querying for vessel")
	}
	v.IP = net.ParseIP(ip)
	return &v, nil
}

// QueryVessels finds all vessels in the database.
func QueryVessels(ctx context.Context, db *sql.DB) ([]*Vessel, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT v.id, v.created_at, INET_NTOA(v.ip), v.cfe, v.data, COALESCE(vpc.clicks, 0)
		FROM vessels v
		LEFT JOIN (
			SELECT vessel_id, COUNT(*) AS clicks
			FROM vessel_pilot_clicks
			GROUP BY vessel_id
		) vpc ON vpc.vessel_id = v.id
		`)
	if err != nil {
		return nil, errors.Wrap(err, "error querying for vessels")
	}
	vessels := []*Vessel{}
	for rows.Next() {
		var (
			ip string
			v  Vessel
		)
		if err := rows.Scan(&v.ID, &v.CreatedAt, &ip, &v.CFE, &v.Data, &v.PilotClicks); err != nil {
			return nil, errors.Wrap(err, "error scanning vessel row")
		}
		v.IP = net.ParseIP(ip)
		vessels = append(vessels, &v)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error scanning vessel rows")
	}
	return vessels, nil
}
