package main

import (
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
			name: "success when content is valid: complete",
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
		{
			name: "success when content is valid: type",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO(dev.lorem): improve package division -> td-design
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Author:      "dev.lorem",
					Type:        "td-design",
					Description: "improve package division",
					Cost:        "",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},
		{
			name: "success when content is valid: description",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO(dev.lorem): improve package division 
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Author:      "dev.lorem",
					Type:        "",
					Description: "improve package division",
					Cost:        "",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},
		{
			name: "success when content is valid: missing author",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO: improve package division 
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Author:      "",
					Type:        "",
					Description: "improve package division",
					Cost:        "",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},
		{
			name: "zero values when SATD is missing",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO: improve package division 
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{},
			err:  nil,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			// <RUN>: output
			got, err := ParseRegex(tt.fileContent, tt.fileName)
			if err != nil {
				t.Errorf("error %v", err)
			}

			if len(got) == 0 {
				t.Fail()
			}

			// <VALIDATE>: verify results
			for _, g := range got {
				for _, w := range tt.want {
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
