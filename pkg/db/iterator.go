package db

type Iterator interface {
	Next() (value interface{}, done bool, err error)
}

func Each(i Iterator, fn func(interface{}) error) error {
	if i != nil {
		for {
			obj, last, err := i.Next()
			if err != nil {
				return err
			}

			if last {
				break
			}

			if err := fn(obj); err != nil {
				return err
			}
		}
	}

	return nil
}
