package moodle_file_renamer

import (
	"fmt"
	"strings"
)

type MoodleFile struct {
	UID         string
	FirstName   string
	LastName    string
	SubmitID    string
	Note        []string
	CurrentName string
}

func buildMoodleFile(entries []string) *MoodleFile {
	mf := &MoodleFile{}
	mf.CurrentName = entries[0]
	mf.UID = entries[1]
	mf.LastName = entries[2]
	mf.FirstName = entries[3]
	mf.SubmitID = entries[4]
	mf.Note = entries[5:]
	return mf
}

func Parse(fileName string) (*MoodleFile, error) {
	entries := strings.Split(fileName, "_")
	if len(entries) < 6 {
		return nil, fmt.Errorf("%s: invalid file name", fileName)
	}
	return &MoodleFile{
		UID:       entries[0],
		LastName:  entries[1],
		FirstName: entries[2],
		SubmitID:  entries[3],
		Note:      entries[4:],
	}, nil
}

func (f *MoodleFile) Format(format string) string {
	result := format
	result = strings.Replace(result, "%uid", f.UID, -1)
	result = strings.Replace(result, "%fname", f.FirstName, -1)
	result = strings.Replace(result, "%lname", f.LastName, -1)
	result = strings.Replace(result, "%sid", f.SubmitID, -1)
	result = strings.Replace(result, "%default", fmt.Sprintf("%s_%s_%s", f.UID, f.LastName, f.FirstName), -1)
	result = strings.Replace(result, "%original", fmt.Sprintf("%s_%s_%s_%s_%s", f.UID, f.LastName, f.FirstName, f.SubmitID, strings.Join(f.Note, "_")), -1)
	result = strings.Replace(result, "%%", "%", -1)
	return result
}
