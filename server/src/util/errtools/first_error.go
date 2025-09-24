package errtools

func GetFirstError(errs <-chan error) error {
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
