package data

import (
	"fmt"
	"runtime"

	//"io"
	"time"

	"github.com/Tokumicn/clickhouse-go/lib/binary"
	"github.com/Tokumicn/clickhouse-go/lib/protocol"
)

type ServerInfo struct {
	Name         string
	Revision     uint64
	MinorVersion uint64
	MajorVersion uint64
	Timezone     *time.Location
}

func (srv *ServerInfo) Read(decoder *binary.Decoder) (err error) {
	if srv.Name, err = decoder.String(); err != nil {
		return fmt.Errorf("could not read server name: %v", err)
	}
	if srv.MajorVersion, err = decoder.Uvarint(); err != nil {
		return fmt.Errorf("could not read server major version: %v", err)
	}
	if srv.MinorVersion, err = decoder.Uvarint(); err != nil {
		return fmt.Errorf("could not read server minor version: %v", err)
	}
	if srv.Revision, err = decoder.Uvarint(); err != nil {
		return fmt.Errorf("could not read server revision: %v", err)
	}
	if srv.Revision >= protocol.DBMS_MIN_REVISION_WITH_SERVER_TIMEZONE {
		timezone, err := decoder.String()
		if err != nil {
			return fmt.Errorf("could not read server timezone: %v", err)
		}
		if srv.Timezone, err = getTimezone(timezone); err != nil {
			return fmt.Errorf("could not load time location: %v", err)
		}
	}
	return nil
}

func getTimezone(timezone string) (*time.Location, error) {
	if runtime.GOOS == "windows" && timezone == "posixrules" {
		timezone = "Local"
	}
	return time.LoadLocation(timezone)
}

func (srv ServerInfo) String() string {
	return fmt.Sprintf("%s %d.%d.%d (%s)", srv.Name, srv.MajorVersion, srv.MinorVersion, srv.Revision, srv.Timezone)
}
