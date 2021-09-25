package manager

import (
	"fmt"
	"io/ioutil"
)

func scan(d string) (map[string][]string, error) {
	var ds = make(map[string][]string)
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if isSystemFile(f.Name())  { // ignore system files
			continue
		}

		if f.IsDir() {
			dFiles, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", d, f.Name()))
			if err != nil {
				return nil, err
			}

			for _, df := range dFiles {
				if !df.IsDir() { // ignore nested directories
					if isSystemFile(df.Name()) { // ignore system files
						continue
					}
					ds[f.Name()] = append(ds[f.Name()], df.Name())
				}
			}
		} else {
			ds[defaultDirName] = append(ds[defaultDirName], f.Name())
		}
	}

	return ds, err
}

func isSystemFile(name string) bool {
	if name[:1] == "." {
		return true
	}
	return false
}