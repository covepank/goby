package bufferpool

var (
	_pool = NewPool()
	// Get 从全局 _pool 中获取 buffer
	Get = _pool.Get
)
