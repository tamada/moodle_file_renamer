package moodle_file_renamer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/tamada/moodle_file_renamer/utils"
)

type MoodleDir interface {
	ReadDir() ([]*MoodleFile, error)
	Close() error
	Restore() error
	Rename(format string) error
}

type moodleDir struct {
	path   string
	mapper map[string]*MoodleFile
}

func (md *moodleDir) Restore() error {
	return md.Rename("%original")
}

func (md *moodleDir) Rename(format string) error {
	mfs, err := md.ReadDir()
	if err != nil {
		return err
	}
	var errs []error
	for _, mf := range mfs {
		newName := mf.Format(format)
		newPath := fmt.Sprintf("%s/%s", md.path, newName)
		if err := os.Rename(fmt.Sprintf("%s/%s", md.path, mf.CurrentName), newPath); err != nil {
			errs = append(errs, err)
		}
		delete(md.mapper, mf.CurrentName)
		md.mapper[newName] = mf
		mf.CurrentName = newName
	}
	return errors.Join(errs...)
}

func (md *moodleDir) Close() error {
	writer, err := os.OpenFile(md.mapperPath(), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	out := csv.NewWriter(writer)
	for current, mf := range md.mapper {
		data := []string{current, mf.UID, mf.LastName, mf.FirstName, mf.SubmitID}
		data = append(data, mf.Note...)
		out.Write(data)
	}
	out.Flush()
	return nil
}

func (md *moodleDir) ReadDir() ([]*MoodleFile, error) {
	var result []*MoodleFile
	entries, err := os.ReadDir(md.path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		name := entry.Name()
		mf, ok := md.mapper[name]
		if !ok {
			mf, err = Parse(name)
			if err != nil {
				continue
			}
			mf.CurrentName = name
			md.mapper[name] = mf
		}
		result = append(result, mf)
	}
	return result, nil
}

func Open(path string) (MoodleDir, error) {
	if !utils.IsDir(path) {
		return nil, fmt.Errorf("%s: not a directory", path)
	}
	md := &moodleDir{path: path}
	err := md.initialize()
	return md, err
}

func (md *moodleDir) initialize() error {
	md.mapper = make(map[string]*MoodleFile)
	if md.existMapperFile() {
		return md.readMapperFile()
	}
	return nil
}

func (md *moodleDir) readMapperFile() error {
	file, err := os.Open(md.mapperPath())
	if err != nil {
		return err
	}
	defer file.Close()
	entries, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		md.mapper[entry[0]] = buildMoodleFile(entry)
	}
	return nil
}

func (md *moodleDir) existMapperFile() bool {
	stat, err := os.Stat(md.mapperPath())
	return err == nil && stat.Mode().IsRegular()
}

func (md *moodleDir) mapperPath() string {
	return fmt.Sprintf("%s/.mfr.csv", md.path)
}
