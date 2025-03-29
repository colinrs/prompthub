package utils

const (
	defaultPage     = 1
	defaultPageSize = 10
)

func Page(page, pageSize int) (int, int) {

	if page < 1 {
		page = defaultPage
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return (page - 1) * pageSize, pageSize
}
