# Todo to Tech Debt (tttd)

The aim of this tool is to help SWE Teams in the task of better controlling it's technical debts by
automatizing the extraction of Self Admited Technical Debts (SATDs) from the code to the used issues manager.

## How it works
    - Detect: search the code, line by line, for a TODO mark with description and/or reporter, type and cost 
    - Transform: transform the SATDs to CSV format (Github/Jira)
    - Export: exports to the default or given path

## Example of valid SATDs
```
func validsSatd() {
	fmt.Println("Hello world!")

	// TODO: simple
	// TODO(simple): with description
	// TODO(simple): with description > type
	// TODO(simple): with description > type => cost

	// TODO: description without user > improve package division.
	fmt.Println("Hello tttd!")

	// TODO(al.lo): td-maintenence > improve package division.
	fmt.Println("Hello tttd!")
}
```

## How to run
$ go run . '["test.go"]'

## How to use
You need to create a workflow file in your project:

your-project/.github/workflows/extract-satds.yaml
```
name: Extract SATDs to CSV

on: [push]

permissions:
  contents: write

jobs:
  extract-satds-to-csv:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Run Technical Debt Analysis 
      uses: feilerdev/tttd@v0.1.5.1
    - name: Commit results
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add .
        git commit -m "Add technical debt analysis results" -a || echo "No changes to commit"
    - name: Push changes
      uses: ad-m/github-push-action@master
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        branch: ${{ github.ref }}
```
