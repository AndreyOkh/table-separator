package ods

type ODS struct {
	Content struct {
		Body struct {
			Spreadsheet struct {
				Table []struct {
					TableRow []struct {
						TableCell []struct {
							P string `xml:"p" json:"p,omitempty"`
						} `xml:"table-cell" json:"table-cell,omitempty"`
					} `xml:"table-row" json:"table-row,omitempty"`
				} `xml:"table" json:"table,omitempty"`
			} `xml:"spreadsheet" json:"spreadsheet,omitempty"`
		} `xml:"body" json:"body,omitempty"`
	} `json:"content"`
}
