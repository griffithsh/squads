package data

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var internalNames = map[string][]string{
	"Alyssa":       []string{"F"},
	"Arneldo":      []string{"M"},
	"Atlus":        []string{"M"},
	"Axis":         []string{"M"},
	"Bacon":        []string{"M"},
	"Balustrade":   []string{"F"},
	"Belle":        []string{"F", "Pure"},
	"Benthamen":    []string{"M"},
	"Bentholemew":  []string{"M"},
	"Bolus":        []string{"M"},
	"Callisto":     []string{"F"},
	"Callo":        []string{"M"},
	"Cristian":     []string{"M"},
	"Devuan":       []string{"M"},
	"Divernon":     []string{"F"},
	"Donkey":       []string{"M"},
	"Dungaree":     []string{"M"},
	"Edward":       []string{"M"},
	"Eloa":         []string{"F"},
	"Euphemia":     []string{"F"},
	"Fankrastha":   []string{"F"},
	"Frederick":    []string{"M"},
	"Gao":          []string{"F"},
	"Gerald":       []string{"M"},
	"Gordania":     []string{"F"},
	"Hana":         []string{"F"},
	"Harmonia":     []string{"F"},
	"Helloise":     []string{"F"},
	"Humperdink":   []string{"M"},
	"Ignatius":     []string{"M"},
	"Ignold":       []string{"M"},
	"Iscariot":     []string{"M", "Villanouse"},
	"Ismalloray":   []string{"F"},
	"Jahnsenn":     []string{"M"},
	"Jamieson":     []string{"M"},
	"Jannifern":    []string{"F"},
	"Jan":          []string{"M", "F"},
	"Kalisto":      []string{"F"},
	"Kamio":        []string{"F"},
	"Katherita":    []string{"F"},
	"Ketlin":       []string{"F"},
	"Krastin":      []string{"M"},
	"Lanneth":      []string{"F"},
	"Legothory":    []string{"F"},
	"Lucas":        []string{"M"},
	"Maillorne":    []string{"F"},
	"Mattieson":    []string{"M"},
	"Nelson":       []string{"M"},
	"Nolan":        []string{"M"},
	"Nostory":      []string{"F"},
	"Ollivene":     []string{"F"},
	"Ormond":       []string{"M"},
	"Oswalt":       []string{"M"},
	"Panseur":      []string{"M"},
	"Pantry":       []string{"M"},
	"Perogue":      []string{"M"},
	"Polter":       []string{"M"},
	"Punt":         []string{"M"},
	"Pursivonian":  []string{"F"},
	"Quincy":       []string{"M"},
	"Qui":          []string{"F"},
	"Ramathese":    []string{"M"},
	"Rederick":     []string{"M"},
	"Rimcy":        []string{"F"},
	"Roon":         []string{"M"},
	"Sallivoce":    []string{"F"},
	"Samithee":     []string{"M"},
	"Satchel":      []string{"M", "F"},
	"Sera":         []string{"F"},
	"Shanto":       []string{"F"},
	"Staunton":     []string{"M"},
	"Thames":       []string{"M"},
	"Theodora":     []string{"F"},
	"Timjamen":     []string{"M"},
	"Undine":       []string{"F"},
	"Unicerve":     []string{"M"},
	"Variose":      []string{"M"},
	"Victohia":     []string{"F"},
	"Violet":       []string{"F"},
	"Volturbulent": []string{"M"},
	"Winchester":   []string{"F"},
	"Xactabol":     []string{"M"},
	"Xerxes":       []string{"M"},
	"Xin":          []string{"F"},
	"Yalladin":     []string{"M"},
	"Yellow":       []string{"F"},
	"Yossarian":    []string{"M"},
	"Zenta":        []string{"F"},
	"Zod":          []string{"M"},
	"Zomparion":    []string{"M"},
}

// Names returns ...
func (a *Archive) Names() map[string][]string {
	return a.names
}

func parseNames(r io.Reader) (map[string][]string, error) {
	result := map[string][]string{}

	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			if line == "" {
				break

			}
		} else if err != nil {
			return nil, fmt.Errorf("read line: %v", err)
		}

		parts := strings.SplitN(line, ":", 2)
		name := strings.Trim(parts[0], " \n\r\t")
		if name == "" {
			continue
		}
		tags := []string{}

		if len(parts) > 1 {
			tags = strings.Split(parts[1], ",")
		}

		cleaned := tags[:0]
		for _, tag := range tags {

			clean := strings.Trim(tag, " \n\r\t")
			if clean == "" {
				continue
			}
			cleaned = append(cleaned, clean)
		}
		result[name] = cleaned
	}

	return result, nil
}
