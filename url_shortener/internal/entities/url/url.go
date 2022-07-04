package url

type URL struct {
	ID      int64  `json:"id"`
	Created string `json:"created"`
	Long    string `json:"long"`
	Short   string `json:"short"`
	Admin   string `json:"admin"`
}
