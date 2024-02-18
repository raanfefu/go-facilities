package common

import "fmt"

const (
	TLS  ModeType = "https"
	HTTP ModeType = "http"
)

type ModeType string

func (c ModeType) String() string {
	return string(c)
}

func ParseMode(s *string) (c *ModeType, err error) {
	if s != nil {
		capabilities := map[ModeType]struct{}{
			HTTP: {},
			TLS:  {},
		}

		mode := ModeType(*s)
		_, ok := capabilities[mode]
		if !ok {
			return c, fmt.Errorf("cannot parse: %s as mode", *s)
		}
		return &mode, nil
	} else {
		return nil, fmt.Errorf("cannot parse: %s  as mode", *s)
	}
}
