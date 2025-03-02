package collection

// CollectAll collects all logs from the configured sources, moving them to the outbox folder.
// Returns the number of files collected.
func CollectAll(sources []Source) (int, error) {
	var collected int
	for _, source := range sources {
		moved, err := source.Collect()
		if err != nil {
			return collected, err
		}
		collected += moved
	}
	return collected, nil
}
