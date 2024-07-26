package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
)

func main() {
	req := &folders.FetchFolderRequest{
		OrgID: uuid.FromStringOrNil(folders.DefaultOrgID),
	}

	res, err := folders.GetAllFolders(req, folders.FetchFolderDependencies{
		DataFetcher: folders.DefaultFetcher{},
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	folders.PrettyPrint(res)
}
