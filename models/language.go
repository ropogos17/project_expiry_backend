package models

// Languages represents multilingual text entries
type Languages struct {
	ID int    `json:"id" db:"id"`
	En string `json:"en" db:"en"` // English
	Hu string `json:"hu" db:"hu"` // Hungarian
	Ro string `json:"ro" db:"ro"` // Romanian
	Es string `json:"es" db:"es"` // Spanish
	Zh string `json:"zh" db:"zh"` // Chinese (Mandarin)
	Fr string `json:"fr" db:"fr"` // French
	De string `json:"de" db:"de"` // German
	Pt string `json:"pt" db:"pt"` // Portuguese
	Ru string `json:"ru" db:"ru"` // Russian
	Ua string `json:"ua" db:"ua"` // Ukrainian
	Ar string `json:"ar" db:"ar"` // Arabic
}
