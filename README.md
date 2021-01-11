# Go Project

An application which scrape from a given URL the HTML version of page, Page Title, Headings, Links (external, internal, inaccessible) and Logins.

## Building locally

1.Clone the project

2.Run `go mod download && go build main`

3.Run `go run main.go`

## Usage

The project is a simple REST application with a GET endpoint which runs locally at port 7171. You can give it a try with Postman.
And as a query parameter `url` paste the link of a page which one you want to scrape. The successful response would be in JSON format e.g.:


    {
    "HtmlVersion": "<e.g HTML Version 5>",
    "PageTitle": " <Page Title> ",
    "Headings": {
        "h1": <amount of h1 elements e.g 1>,
        "h2": <e.g 3>,
        "h3": <e.g 3>,
        "h4": <e.g 16>,
        "h5": <e.g 1>
    },
    "Links": {
        "external": <amount of external links e.g. 22>,
        "inaccessible": <amount of inaccessible links e.g 2>,
        "internal": <amount of internal links e.g. 4>
    },
    "Login": <amount of logins e.g. 1>
    }


