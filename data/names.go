package data

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var internalNames = map[string][]string{
	"Alyssa":       {"F"},
	"Arneldo":      {"M"},
	"Atlus":        {"M"},
	"Axis":         {"M"},
	"Bacon":        {"M"},
	"Balustrade":   {"F"},
	"Belle":        {"F", "Pure"},
	"Benthamen":    {"M"},
	"Bentholemew":  {"M"},
	"Bolus":        {"M"},
	"Callisto":     {"F"},
	"Callo":        {"M"},
	"Cristian":     {"M"},
	"Devuan":       {"M"},
	"Divernon":     {"F"},
	"Donkey":       {"M"},
	"Dungaree":     {"M"},
	"Edward":       {"M"},
	"Eloa":         {"F"},
	"Euphemia":     {"F"},
	"Fankrastha":   {"F"},
	"Frederick":    {"M"},
	"Gao":          {"F"},
	"Gerald":       {"M"},
	"Gordania":     {"F"},
	"Hana":         {"F"},
	"Harmonia":     {"F"},
	"Helloise":     {"F"},
	"Humperdink":   {"M"},
	"Ignatius":     {"M"},
	"Ignold":       {"M"},
	"Iscariot":     {"M", "Villanouse"},
	"Ismalloray":   {"F"},
	"Jahnsenn":     {"M"},
	"Jamieson":     {"M"},
	"Jannifern":    {"F"},
	"Jan":          {"M", "F"},
	"Kalisto":      {"F"},
	"Kamio":        {"F"},
	"Katherita":    {"F"},
	"Ketlin":       {"F"},
	"Krastin":      {"M"},
	"Lanneth":      {"F"},
	"Legothory":    {"F"},
	"Lucas":        {"M"},
	"Maillorne":    {"F"},
	"Mattieson":    {"M"},
	"Nelson":       {"M"},
	"Nolan":        {"M"},
	"Nostory":      {"F"},
	"Ollivene":     {"F"},
	"Ormond":       {"M"},
	"Oswalt":       {"M"},
	"Panseur":      {"M"},
	"Pantry":       {"M"},
	"Perogue":      {"M"},
	"Polter":       {"M"},
	"Punt":         {"M"},
	"Pursivonian":  {"F"},
	"Quincy":       {"M"},
	"Qui":          {"F"},
	"Ramathese":    {"M"},
	"Rederick":     {"M"},
	"Rimcy":        {"F"},
	"Roon":         {"M"},
	"Sallivoce":    {"F"},
	"Samithee":     {"M"},
	"Satchel":      {"M", "F"},
	"Sera":         {"F"},
	"Shanto":       {"F"},
	"Staunton":     {"M"},
	"Thames":       {"M"},
	"Theodora":     {"F"},
	"Timjamen":     {"M"},
	"Undine":       {"F"},
	"Unicerve":     {"M"},
	"Variose":      {"M"},
	"Victohia":     {"F"},
	"Violet":       {"F"},
	"Volturbulent": {"M"},
	"Winchester":   {"F"},
	"Xactabol":     {"M"},
	"Xerxes":       {"M"},
	"Xin":          {"F"},
	"Yalladin":     {"M"},
	"Yellow":       {"F"},
	"Yossarian":    {"M"},
	"Zenta":        {"F"},
	"Zod":          {"M"},
	"Zomparion":    {"M"},
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
