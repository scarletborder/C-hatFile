package search

// type SearchRequest struct{

// 	Tags string ``
// }

type SearchResult struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []string `json:"tags"`
}
