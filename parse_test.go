package main

import (
	"fmt"
	"testing"
)

type TestCase struct {
	name        string
	fileContent string
	fileName    string
	want        []TechnicalDebt
	err         error
}

func Test_ParseRegex(t *testing.T) {
	t.Parallel()

	testCases := []TestCase{
		{
			name: "success when content is valid",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO(dev.lorem): improve package division -> td-design => $$
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Author:      "dev.lorem",
					Type:        "td-design",
					Description: "improve package division",
					Cost:        "$$",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// <RUN>: output
			got, err := ParseRegex(tt.fileContent, tt.fileName)
			if err != nil {
				t.Errorf("error %v", err)
			}

			// <VALIDATE>: verify results
			for _, g := range got {
				for _, w := range tt.want {
					fmt.Println(g.Author)
					fmt.Println(g.Description)
					fmt.Println(g.Type)
					fmt.Println(g.Cost)
					fmt.Println(g.File)
					fmt.Println(g.Line)

					if g.Author != w.Author {
						t.Errorf("author: got %q, wanted %q", g.Author, w.Author)
					}

					if g.Type != w.Type {
						t.Errorf("type: got %q, wanted %q", g.Type, w.Type)
					}

					if g.Description != w.Description {
						t.Errorf("description: got %q, wanted %q", g.Description, w.Description)
					}

					if g.Cost != w.Cost {
						t.Errorf("cost: got %q, wanted %q", g.Cost, w.Cost)
					}

					if g.File != w.File {
						t.Errorf("file: got %q, wanted %q", g.File, w.File)
					}

					if g.Line != w.Line {
						t.Errorf("line: got %q, wanted %q", g.Line, w.Line)
					}
				}
			}
		})
	}
}
