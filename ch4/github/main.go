package main

import (
	"encoding/json"
	"fmt"
	htmlTemplate "html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResults struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

func (is Issue) String() string {
	return fmt.Sprintf("Created at: %s\nTitle: %q\nState: %q\n", is.CreatedAt.Format("01/02/2006 15:04:05"), is.Title, is.State)
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResults, error) {
	q := url.QueryEscape(strings.Join(terms, " "))

	res, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", res.Status)
	}

	var results IssuesSearchResults
	if err := json.NewDecoder(res.Body).Decode(&results); err != nil {
		res.Body.Close()
		return nil, err
	}

	res.Body.Close()
	return &results, nil
}

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

const temp = `{{.TotalCount}} issues:
{{range .Items}}-----------------------------
Number:{{.Number}}
User:{{.User.Login}}
Title:{{.Title | printf "%.64s"}}
Age:{{.CreatedAt | daysAgo}} days
{{end}}`

var issuesList = htmlTemplate.Must(htmlTemplate.New("issuelist").Parse(`
<H1>ISSUES</H1>
<h3>{{.TotalCount}}</h3>
<table>
<tr style='text-align: left;'>
	<th>State</th>
	<th>User</th>
	<th>Title</th>
</tr>
{{range .Items}}
<tr>
	<td>{{.State}}</td>
	<td><a href="{{.User.HTMLURL}}">{{.User.Login}}</a></td>
	<td><a href="{{.HTMLURL}}">{{.Title}}</a></td>
</tr>
{{end}}
</table>
`))

func main() {
	issues, err := SearchIssues([]string{"is:open", "repo:golang/go 3133 10535"})
	if err != nil {
		log.Fatalln(err)
	}
	// printIssuesByDateRange(issues)

	// report, err := template.New("report").Funcs(template.FuncMap{"daysAgo": daysAgo}).Parse(temp)
	// if err != nil {
	// log.Fatalln(err)
	// }

	if err := issuesList.Execute(os.Stdout, issues); err != nil {
		log.Fatalln(err)
	}

}

func printIssuesByDateRange(issues *IssuesSearchResults) {
	issYoungerThanMonth := make([]Issue, 0, issues.TotalCount/3)
	issYoungerThanYear := make([]Issue, 0, issues.TotalCount/3)
	issOlder := make([]Issue, 0, issues.TotalCount/3)

	t := time.Now()
	tMonthBack := time.Date(t.Year(), t.Month()-1, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	tYearBack := time.Date(t.Year()-1, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())

	fmt.Println(t.Format("01/02/2006 15:04:05"))
	fmt.Println(tMonthBack.Format("01/02/2006 15:04:05"))
	fmt.Println(tYearBack.Format("01/02/2006 15:04:05"))

	for _, v := range issues.Items {

		if v.CreatedAt.After(tMonthBack) {
			issYoungerThanMonth = append(issYoungerThanMonth, *v)
			continue
		}
		if v.CreatedAt.After(tYearBack) {
			issYoungerThanYear = append(issYoungerThanYear, *v)
			continue
		}
		issOlder = append(issOlder, *v)
	}

	fmt.Printf("\n\nIssues not older than a month:\n%s", issYoungerThanMonth)
	fmt.Printf("\n\nIssues not older than a year:\n%s", issYoungerThanYear)
	fmt.Printf("\n\nIssues at lease one year old:\n%s", issOlder)
}
