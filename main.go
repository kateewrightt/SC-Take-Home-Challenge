package main

import (
	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
)

func main() {
	req := &folders.FetchFolderRequest{
		OrgID: uuid.FromStringOrNil(folders.DefaultOrgID),
	}

	// Initialise the dependencies for fetching folders - using local sample data here
	// We use dependency injection to make the code more testable and modular
	deps := folders.FetchFolderDependencies{
		DataFetcher: folders.DefaultFetcher{},
	}

	// Ideally these variables should not be hardcoded
	pageSize := 2         // Change based on desired page size for pagination
	usePagination := true // Set this to false if you want to use GetAllFolders

	if usePagination {
		folders.FetchAndPrintFoldersWithPagination(req, deps, pageSize)
	} else {
		folders.FetchAndPrintAllFolders(req, deps)
	}
}
