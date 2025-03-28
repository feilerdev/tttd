package main

import (
	"testing"
)

type TestCase struct {
	name        string
	fileContent string
	filePath    string
	fileName    string
	want        []TechnicalDebt
	err         error
}

func Test_ParseRegex(t *testing.T) {
	t.Parallel()

	testCases := []TestCase{
		{
			name:     "success when content is valid: complete",
			filePath: "./tests/",
			fileName: "test.go",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO(dev.lorem): T1 improve package division -> td-design => $$
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Author:      "dev.lorem",
					Type:        "td-design",
					Description: "T1 improve package division",
					Cost:        "$$",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},
		{
			name:     "success when content is valid: type",
			filePath: "./tests/",
			fileName: "test.go",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO(dev.lorem): T2 improve package division -> td-design
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Author:      "dev.lorem",
					Type:        "td-design",
					Description: "T2 improve package division",
					Cost:        "",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},

		{
			name:     "success when no content in a comment",
			filePath: "./tests/",
			fileName: "test.go",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// T3 testing a comment
				fmt.Println("Hello tttd!")
			}`,
			want: nil,
			err:  nil,
		},
		{
			name:     "success when comment is empty",
			filePath: "./tests/",
			fileName: "test.go",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// 
				fmt.Println("Hello tttd!")
			}`,
			want: nil,
			err:  nil,
		},
		{
			name:     "success when content is valid: description",
			filePath: "./tests/",
			fileName: "test.go",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO(dev.lorem): T4 improve package division 
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Author:      "dev.lorem",
					Type:        "",
					Description: "T4 improve package division",
					Cost:        "",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},
		{
			name:     "success when content is valid: missing author",
			filePath: "./tests/",
			fileName: "test.go",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO: T5 improve package division 
				fmt.Println("Hello tttd!")
			}`,
			want: []TechnicalDebt{
				{
					Author:      "",
					Type:        "",
					Description: "T5 improve package division",
					Cost:        "",
					File:        "",
					Line:        8,
				},
			},
			err: nil,
		},
		{
			name:     "zero values when SATD is missing",
			filePath: "./tests/",
			fileName: "test.go",
			fileContent: `package main

			import "fmt"

			func main() {
				fmt.Println("Hello world!")

				// TODO: T6 improve package division 
				fmt.Println("Hello tttd!")
			}`,
			want: nil,
			err:  nil,
		},
	}

	for _, tt := range testCases {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			// <RUN>: output
			got, err := ParseRegex(tt.fileContent, tt.filePath, tt.fileName)
			if err != nil {
				t.Errorf("error %v", err)
			}

			if len(got) != 0 {
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
			}
		})
	}
}
