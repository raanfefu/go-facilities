package configure

import (
	"flag"
	"fmt"
)

func (d *DefaultConfiguraionService) JsonVar(p *interface{}, name string, usage string) {
	if d.jsonFiles == nil {
		d.jsonFiles = make(map[string]*JsonFileArgsType)
	}
	crt := JsonFileArgsType{
		Content: p,
	}
	flag.StringVar(&crt.Name, fmt.Sprintf("%s%s%s%sjson", *d.key, SEPARATOR, name, SEPARATOR), "", usage)
	d.jsonFiles[fmt.Sprintf("%s:%s", *d.key, name)] = &crt

}

func (d *DefaultConfiguraionService) GetJsonNames() []string {

	if d.jsonFiles != nil {
		keys := make([]string, 0)
		for k := range d.jsonFiles {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

func (d *DefaultConfiguraionService) GetJsonValue(name string) *JsonFileArgsType {
	if d.jsonFiles != nil {
		return d.jsonFiles[fmt.Sprintf("%s:%s", *d.key, name)]
	}
	return nil
}
