package geoip

import (
	"log/slog"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

type dbRecord struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

// DB is a GeoIP database provider
type DB struct {
	reader *maxminddb.Reader
}

// New GeoIP database provider
func New(filename string) (*DB, error) {
	if filename == "" {
		return nil, nil
	}
	reader, err := maxminddb.Open(filename)
	if err != nil {
		return nil, err
	}
	slog.Debug("using geo IP database", "uri", filename)
	return &DB{
		reader: reader,
	}, nil
}

// Close GeoIP database
func (db *DB) Close() error {
	if db.reader != nil {
		return db.reader.Close()
	}
	return nil
}

// LookupCountry find IP country code
func (db *DB) LookupCountry(ip net.IP) (string, error) {
	record := &dbRecord{}
	err := db.reader.Lookup(ip, &record)
	if err != nil {
		return "", err
	}
	return record.Country.ISOCode, nil
}
