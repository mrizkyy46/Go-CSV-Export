package model

type Proxy struct {
	Amt   string
	Descr string
	Date  string
	ID    string
}

type Source struct {
	Date        string
	ID          string
	Amount      string
	Description string
}

type Export struct {
	Amt     string
	Descr   string
	Date    string
	ID      string
	Remarks string
}
