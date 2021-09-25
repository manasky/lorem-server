package manager

import (
	"math/rand"
	"strings"
	"time"
)

const (
	defaultDirName = "__root"
	CacheDir = ".cache"
)

type Manager struct {
	base string
	ds   map[string][]string
	mc   []string // index of categories
}

func New(dir string) (*Manager, error) {
	ds, err := scan(dir)
	if err != nil {
		return nil, err
	}

	var mc []string
	for c, _ := range ds {
		mc = append(mc, c)
	}

	return &Manager{
		base: dir,
		ds:   ds,
		mc:   mc,
	}, nil
}

func (m *Manager) Pick(cat string) string {
	if cat == "" {
		cat = m.randomCategory()
	}
	return m.randomEntity(cat)
}

func (m *Manager) Total() int {
	var t int
	for _, d := range m.ds {
		t += len(d)
	}
	return t
}

func (m *Manager) randomCategory() string {
	rand.Seed(time.Now().UnixNano())
	return m.mc[rand.Intn(len(m.mc))]
}

func (m *Manager) randomEntity(c string) string {
	rand.Seed(time.Now().UnixNano())
	if _, ok := m.ds[c]; !ok {
		return ""
	}

	names := []string{m.base}
	if c != defaultDirName {
		names = append(names, c)
	}
	names = append(names, m.ds[c][rand.Intn(len(m.ds[c]))])

	return strings.Join(names, "/")
}
