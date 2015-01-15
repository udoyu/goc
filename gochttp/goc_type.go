package gochttp

type str_t struct {
	data   *string
	result int
}

type chan_data_t struct {
	data       []chan *str_t
	is_disable bool
}

func NewChanStr(size int) *chan_data_t {
	t := new(chan_data_t)
	t.data = make([]chan *str_t, size)
	for i := 0; i < size; i++ {
		t.data[i] = make(chan *str_t, 1)
	}
	return t
}

func (p *chan_data_t) Add(data *str_t, i int) {
	if !p.is_disable {
		p.data[i] <- data
	}
}

func (p *chan_data_t) Get(i int) *str_t {
	return <-p.data[i]
}
