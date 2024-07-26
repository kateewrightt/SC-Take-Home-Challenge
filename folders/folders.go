package folders

import (
	"github.com/gofrs/uuid"
)

// GetAllFolders fetches all folders for an organisation
// This function now accepts dependencies to facilitate dependency injection
func GetAllFolders(req *FetchFolderRequest, deps FetchFolderDependencies) (*FetchFolderResponse, error) {

	// Unused Variables
	// var (
	//	err error
	//	f1  Folder
	//	fs  []*Folder
	// )
	// f := []Folder{}

	// We now use the DataFetcher to fetch all folders for an Organisation ID
	r, err := FetchAllFoldersByOrgID(req.OrgID, deps.DataFetcher)

	// Add error handling
	if err != nil {
		return nil, err
	}

	// Converts map vaules to a slice of folders
	// This loop iterates over the map returned by FetchAllFoldersByOrgID and converts
	// each value to a folder struct and then appends each folder struct to the slice f
	// for k, v := range r {
	//	f = append(f, *v)
	// }

	// Converts the slice of folders to a slice of pointers to folder
	// This second loop then takes the slice of folder structs f and
	// converts it back to a slice of pointers to folder struct fp
	//var fp []*Folder
	//for k1, v1 := range f {
	//	fp = append(fp, &v1)
	//}

	// Uses previous code here that is unecessary
	//var ffr *FetchFolderResponse
	//ffr = &FetchFolderResponse{Folders: fp}
	//return ffr, nil

	// Simplified creation of FetchFolderResponse
	// We directly assign the folders returned by FetchAllFoldersByOrgID to the response
	ffr := &FetchFolderResponse{Folders: r}
	return ffr, nil
}

// Interface for loading folders data in memory
// Worth noting loading data in memory and query/filter in memory
// Is always a bad idea in most cases as opposed to leveraging a DB solution and their querying
type DataFetcherInterface interface {
	GetFolders() []*Folder
}

// DefaultFetcher is a default implementation of DataFetcherInterface
// We use this to fetch Sample data from a local source
type DefaultFetcher struct{}

// GetFolders retrieves sample data (default implementation)
func (f DefaultFetcher) GetFolders() []*Folder {
	return GetSampleData()
}

// FetchAllFoldersByOrgID retrieves all folders for a given organisation ID
// We use the DataFetcherInterface to obtain data
func FetchAllFoldersByOrgID(orgID uuid.UUID, dataFetcher DataFetcherInterface) ([]*Folder, error) {
	// We want the retrieval of data here to be a modular component in our app
	// Here the default method is getting data through local json
	// We want to be able to generate this data in tests too
	// To faciliate this we will change the hardcoded function call into a interface
	// And inject the data retrieval as an application dependency
	folders := dataFetcher.GetFolders()

	resFolder := []*Folder{}
	for _, folder := range folders {
		if folder.OrgId == orgID {
			resFolder = append(resFolder, folder)
		}
	}
	return resFolder, nil
}
