package storage

type StorageFilters struct {
	names []string
}

type StorageFilter func(o *StorageFilters)

func FilterNames(names []string) StorageFilter {
	return func(f *StorageFilters) {
		f.names = names
	}
}

func FilterName(name string) StorageFilter {
	return func(f *StorageFilters) {
		f.names = append(f.names, name)
	}
}
