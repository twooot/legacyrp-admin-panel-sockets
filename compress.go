package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

type CCharacter struct {
	Dead     bool   `json:"a,omitempty"`
	FullName string `json:"b,omitempty"`
	ID       int64  `json:"c,omitempty"`
	InShell  bool   `json:"d,omitempty"`
}

type CVehicle struct {
	Driving bool   `json:"a,omitempty"`
	ID      int64  `json:"b,omitempty"`
	Model   string `json:"c,omitempty"`
	Name    string `json:"d,omitempty"`
}

type CPlayer struct {
	AFK            int64       `json:"a,omitempty"`
	Character      *CCharacter `json:"b,omitempty"`
	Movement       string      `json:"c,omitempty"`
	Invisible      bool        `json:"d,omitempty"`
	InvisibleSince int64       `json:"e,omitempty"`
	Name           string      `json:"f,omitempty"`
	Source         int64       `json:"g,omitempty"`
	Steam          string      `json:"h,omitempty"`
	Vehicle        *CVehicle   `json:"i,omitempty"`
}

func CompressPlayers(players []map[string]interface{}) []CPlayer {
	compressed := make([]CPlayer, len(players))

	for i, p := range players {
		character := getMap("character", p)
		var c *CCharacter
		if character != nil {
			c = &CCharacter{
				Dead:     getBool("dead", character),
				FullName: getString("fullName", character),
				ID:       getInt64("id", character),
				InShell:  getBool("inShell", character),
			}
		}

		vehicle := getMap("vehicle", p)
		var v *CVehicle
		if vehicle != nil {
			v = &CVehicle{
				Driving: getBool("driving", vehicle),
				ID:      getInt64("id", vehicle),
				Model:   getString("model", vehicle),
				Name:    getString("name", vehicle),
			}
		}

		compressed[i] = CPlayer{
			AFK:            getInt64("afk", p),
			Character:      c,
			Movement:       getMovementData(p),
			Invisible:      getBool("invisible", p),
			InvisibleSince: getInt64("invisible_since", p),
			Name:           getString("name", p),
			Source:         getInt64("source", p),
			Steam:          getString("steamIdentifier", p),
			Vehicle:        v,
		}
	}

	return compressed
}

func getMovementData(m map[string]interface{}) string {
	c := getMap("coords", m)

	if c != nil {
		x, xOk := c["x"].(float64)
		y, yOk := c["y"].(float64)
		z, zOk := c["z"].(float64)

		h := getFloat64("heading", m)
		s := getFloat64("speed", m)

		if xOk && yOk && zOk {
			str := fmt.Sprintf("%.1f,%.1f,%.1f,%.1f", x, y, z, h)

			if s != 0 {
				str += fmt.Sprintf(",%.1f", s)
			}

			return str
		}
	}

	return ""
}

func getFloat64(key string, m map[string]interface{}) float64 {
	v, ok := m[key]

	if ok && v != nil {
		f, ok := v.(float64)

		if ok {
			return f
		}
	}

	return 0
}

func getInt64(key string, m map[string]interface{}) int64 {
	v, ok := m[key]

	if ok && v != nil {
		i, ok := v.(int64)

		if !ok {
			f, ok2 := v.(float64)
			if ok2 {
				ok = true
				i = int64(f)
			}
		}

		if ok {
			return i
		}
	}

	return 0
}

func getString(key string, m map[string]interface{}) string {
	v, ok := m[key]

	if ok && v != nil {
		s, ok := v.(string)

		if ok {
			return s
		}
	}

	return ""
}

func getBool(key string, m map[string]interface{}) bool {
	v, ok := m[key]

	if ok && v != nil {
		b, ok := v.(bool)

		if ok {
			return b
		}
	}

	return false
}

func getMap(key string, m map[string]interface{}) map[string]interface{} {
	v, ok := m[key]

	if ok && v != nil {
		m, ok := v.(map[string]interface{})

		if ok {
			return m
		}
	}

	return nil
}

func gzipBytes(b []byte) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)

	_, err := w.Write(b)
	if err != nil {
		log.Error("GZIP write failed: " + err.Error())
	}

	err = w.Close()
	if err != nil {
		log.Error("GZIP close failed: " + err.Error())
	}

	return buf.Bytes()
}