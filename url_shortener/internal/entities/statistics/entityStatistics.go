package statistics

type Statistics struct {
	IP      string `json:"ip"`
	Visited string `json:"visited"`
	Count   int    `json:"count"`
	Long    string `json:"long"`
	Short   string `json:"short"`
	Admin   string `json:"admin"`
}
