package utils

import (
	"reflect"
	"testing"
)

func TestFolderFile(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "onefile",
			args: args{
				fileName: "test1234",
			},
			want: map[string]interface{}{
				"": "test1234",
			},
		},
		{
			name: "folder",
			args: args{
				fileName: "folder1/test1234",
			},
			want: map[string]interface{}{
				"":         "folder1/",
				"folder1/": "test1234",
			},
		},
		{
			name: "folder 2",
			args: args{
				fileName: "folder1/folder2/test1234",
			},
			want: map[string]interface{}{
				"":                 "folder1/",
				"folder1/":         "folder2/",
				"folder1/folder2/": "test1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FolderFile(tt.args.fileName)
			result := make(map[string]interface{}, len(got))
			for i := range got {
				result[got[i].Folder] = got[i].Item
			}

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("FolderFile() = %v, want %v", result, tt.want)
			}
		})
	}
}
