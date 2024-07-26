package folders

import "github.com/gofrs/uuid"

type FetchFolderRequest struct {
	OrgID uuid.UUID
}

type FetchFolderDependencies struct {
	DataFetcher DataFetcherInterface
}

type FetchFolderResponse struct {
	Folders []*Folder
}
