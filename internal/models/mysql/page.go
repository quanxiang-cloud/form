package mysql

func pages(page, size int) (int, int) {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 999
	}
	return page, size
}
