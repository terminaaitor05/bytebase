package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	api "github.com/bytebase/bytebase/backend/legacyapi"
	vcsPlugin "github.com/bytebase/bytebase/backend/plugin/vcs"
	"github.com/bytebase/bytebase/backend/store"
	"github.com/bytebase/bytebase/backend/utils"
	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
)

// SheetService implements the sheet service.
type SheetService struct {
	v1pb.UnimplementedSheetServiceServer
	store *store.Store
}

// NewSheetService creates a new SheetService.
func NewSheetService(store *store.Store) *SheetService {
	return &SheetService{
		store: store,
	}
}

// CreateSheet creates a new sheet.
func (s *SheetService) CreateSheet(ctx context.Context, request *v1pb.CreateSheetRequest) (*v1pb.Sheet, error) {
	if request.Sheet == nil {
		return nil, status.Errorf(codes.InvalidArgument, "sheet must be set")
	}
	currentPrincipalID := ctx.Value(common.PrincipalIDContextKey).(int)

	projectResourceID, err := getProjectID(request.Parent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{
		ResourceID: &projectResourceID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get project with resource id %q, err: %s", projectResourceID, err.Error()))
	}
	if project == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q not found", projectResourceID))
	}
	if project.Deleted {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q had deleted", projectResourceID))
	}

	var databaseUID *int
	if request.Sheet.Database != "" {
		instanceResourceID, databaseName, err := getInstanceDatabaseID(request.Sheet.Database)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		instance, err := s.store.GetInstanceV2(ctx, &store.FindInstanceMessage{
			ResourceID: &instanceResourceID,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get instance with resource id %q, err: %s", instanceResourceID, err.Error()))
		}
		if instance == nil {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("instance with resource id %q not found", instanceResourceID))
		}

		database, err := s.store.GetDatabaseV2(ctx, &store.FindDatabaseMessage{
			ProjectID:    &projectResourceID,
			InstanceID:   &instanceResourceID,
			DatabaseName: &databaseName,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get database with name %q, err: %s", databaseName, err.Error()))
		}
		if database == nil {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("database with name %q not found in project %q instance %q", databaseName, projectResourceID, instanceResourceID))
		}
		databaseUID = &database.UID
	}
	storeSheetCreate, err := convertToStoreSheetMessage(project.UID, databaseUID, currentPrincipalID, request.Sheet)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("failed to convert sheet: %v", err))
	}
	sheet, err := s.store.CreateSheetV2(ctx, storeSheetCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to create sheet: %v", err))
	}
	v1pbSheet, err := s.convertToAPISheetMessage(ctx, sheet)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to convert sheet: %v", err))
	}
	return v1pbSheet, nil
}

// GetSheet returns the requested sheet, cutoff the content if the content is too long and the `raw` flag in request is false.
func (s *SheetService) GetSheet(ctx context.Context, request *v1pb.GetSheetRequest) (*v1pb.Sheet, error) {
	projectResourceID, sheetID, err := getProjectResourceIDSheetID(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	sheetIntID, err := strconv.Atoi(sheetID)
	if err != nil || sheetIntID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid sheet id %s, must be positive integer", sheetID))
	}

	project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{
		ResourceID: &projectResourceID,
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %s not found", projectResourceID))
	}
	if project == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %s not found", projectResourceID))
	}
	if project.Deleted {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q had deleted", projectResourceID))
	}

	currentPrincipalID := ctx.Value(common.PrincipalIDContextKey).(int)
	sheet, err := s.store.GetSheetV2(ctx, &api.SheetFind{
		ID:        &sheetIntID,
		LoadFull:  request.Raw,
		ProjectID: &project.UID,
	}, currentPrincipalID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get sheet: %v", err))
	}
	if sheet == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("sheet with id %d not found", sheetIntID))
	}

	v1pbSheet, err := s.convertToAPISheetMessage(ctx, sheet)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to convert sheet: %v", err))
	}
	return v1pbSheet, nil
}

// SearchSheets returns a list of sheets based on the search filters.
func (s *SheetService) SearchSheets(ctx context.Context, request *v1pb.SearchSheetsRequest) (*v1pb.SearchSheetsResponse, error) {
	projectResourceID, err := getProjectID(request.Parent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	currentPrincipalID := ctx.Value(common.PrincipalIDContextKey).(int)

	sheetFind := &api.SheetFind{}
	if projectResourceID != "-" {
		project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{
			ResourceID: &projectResourceID,
		})
		if err != nil {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %s not found", projectResourceID))
		}
		if project.Deleted {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q had deleted", projectResourceID))
		}
		sheetFind.ProjectID = &project.UID
	}

	// TODO(zp): It is difficult to find all the sheets visible to a principal atomically
	// without adding a new store layer method, which has two parts:
	// 1. creator = principal && visibility in (PROJECT, PUBLIC, PRIVATE)
	// 2. creator ! = principal && visibility in (PROJECT, PUBLIC)
	// So we don't allow empty filter for now.
	if request.Filter == "" {
		return nil, status.Errorf(codes.InvalidArgument, "filter should not be empty")
	}

	specs, err := parseFilter(request.Filter)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	for _, spec := range specs {
		switch spec.key {
		case "creator":
			creatorEmail := strings.TrimPrefix(spec.value, "users/")
			if creatorEmail == "" {
				return nil, status.Errorf(codes.InvalidArgument, "invalid empty creator identifier")
			}
			user, err := s.store.GetUser(ctx, &store.FindUserMessage{
				Email: &creatorEmail,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get user: %s", err.Error()))
			}
			if user == nil {
				return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user with email %s not found", creatorEmail))
			}
			switch spec.operator {
			case comparatorTypeEqual:
				sheetFind.CreatorID = &user.ID
				sheetFind.Visibilities = []api.SheetVisibility{api.ProjectSheet, api.PublicSheet, api.PrivateSheet}
			case comparatorTypeNotEqual:
				sheetFind.ExcludedCreatorID = &user.ID
				sheetFind.Visibilities = []api.SheetVisibility{api.ProjectSheet, api.PublicSheet}
				sheetFind.PrincipalID = &user.ID
			default:
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid operator %q for creator", spec.operator))
			}
		case "starred":
			if spec.operator != comparatorTypeEqual {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid operator %q for starred", spec.operator))
			}
			switch spec.value {
			case "true":
				sheetFind.OrganizerPrincipalIDStarred = &currentPrincipalID
			case "false":
				sheetFind.OrganizerPrincipalIDNotStarred = &currentPrincipalID
			default:
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid value %q for starred", spec.value))
			}
		default:
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid filter key %q", spec.key))
		}
	}
	sheetList, err := s.store.ListSheetsV2(ctx, sheetFind, currentPrincipalID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to list sheets: %v", err))
	}

	var v1pbSheets []*v1pb.Sheet
	for _, sheet := range sheetList {
		v1pbSheet, err := s.convertToAPISheetMessage(ctx, sheet)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to convert sheet: %v", err))
		}
		v1pbSheets = append(v1pbSheets, v1pbSheet)
	}
	return &v1pb.SearchSheetsResponse{
		Sheets: v1pbSheets,
	}, nil
}

// UpdateSheet updates a sheet.
func (s *SheetService) UpdateSheet(ctx context.Context, request *v1pb.UpdateSheetRequest) (*v1pb.Sheet, error) {
	if request.Sheet == nil {
		return nil, status.Errorf(codes.InvalidArgument, "sheet cannot be empty")
	}
	if request.UpdateMask == nil {
		return nil, status.Errorf(codes.InvalidArgument, "update mask cannot be empty")
	}
	if request.Sheet.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "sheet name cannot be empty")
	}

	projectResourceID, sheetID, err := getProjectResourceIDSheetID(request.Sheet.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	sheetIntID, err := strconv.Atoi(sheetID)
	if err != nil || sheetIntID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid sheet id %s, must be positive integer", sheetID))
	}

	project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{
		ResourceID: &projectResourceID,
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %s not found", projectResourceID))
	}
	if project == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %s not found", projectResourceID))
	}
	if project.Deleted {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q had deleted", projectResourceID))
	}

	currentPrincipalID := ctx.Value(common.PrincipalIDContextKey).(int)
	sheet, err := s.store.GetSheetV2(ctx, &api.SheetFind{
		ID:        &sheetIntID,
		ProjectID: &project.UID,
	}, currentPrincipalID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get sheet: %v", err))
	}
	if sheet == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("sheet with id %d not found", sheetIntID))
	}

	sheetPatch := &store.PatchSheetMessage{
		ID:        sheet.UID,
		UpdaterID: currentPrincipalID,
	}

	for _, path := range request.UpdateMask.Paths {
		switch path {
		case "title":
			sheetPatch.Name = &request.Sheet.Title
		case "content":
			statement := string(request.Sheet.Content)
			sheetPatch.Statement = &statement
		case "starred":
			if _, err := s.store.UpsertSheetOrganizer(ctx, &api.SheetOrganizerUpsert{
				SheetID:     sheet.UID,
				PrincipalID: currentPrincipalID,
				Starred:     request.Sheet.Starred,
			}); err != nil {
				return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to update sheet organizer: %v", err))
			}
		case "visibility":
			visibility, err := convertToLegacyAPISheetVisibility(request.Sheet.Visibility)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid visibility %q", request.Sheet.Visibility))
			}
			stringVisibility := string(visibility)
			sheetPatch.Visibility = &stringVisibility
		case "payload":
			sheetPatch.Name = &request.Sheet.Payload
		default:
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid update mask path %q", path))
		}
	}
	storeSheet, err := s.store.PatchSheetV2(ctx, sheetPatch)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to update sheet: %v", err))
	}
	v1pbSheet, err := s.convertToAPISheetMessage(ctx, storeSheet)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to convert sheet: %v", err))
	}

	return v1pbSheet, nil
}

// DeleteSheet deletes a sheet.
func (s *SheetService) DeleteSheet(ctx context.Context, request *v1pb.DeleteSheetRequest) (*emptypb.Empty, error) {
	projectResourceID, sheetID, err := getProjectResourceIDSheetID(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	sheetIDInt, err := strconv.Atoi(sheetID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid sheet id %s", sheetID))
	}
	project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{
		ResourceID: &projectResourceID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get project with resource id %q, err: %s", projectResourceID, err.Error()))
	}
	if project == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q not found", projectResourceID))
	}
	if project.Deleted {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q had deleted", projectResourceID))
	}

	currentPrincipalID := ctx.Value(common.PrincipalIDContextKey).(int)

	sheet, err := s.store.GetSheetV2(ctx, &api.SheetFind{
		ID:        &sheetIDInt,
		ProjectID: &project.UID,
	}, currentPrincipalID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get sheet: %v", err))
	}
	if sheet == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("sheet with id %d not found", sheetIDInt))
	}

	if err := s.store.DeleteSheet(ctx, &api.SheetDelete{
		ID:        sheetIDInt,
		DeleterID: currentPrincipalID,
	}); err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to delete sheet: %v", err))
	}

	return &emptypb.Empty{}, nil
}

// SyncSheets syncs sheets from VCS.
func (s *SheetService) SyncSheets(ctx context.Context, request *v1pb.SyncSheetsRequest) (*emptypb.Empty, error) {
	// TODO(tianzhou): uncomment this after adding the test harness to using Enterprise version.
	// if !s.licenseService.IsFeatureEnabled(api.FeatureVCSSheetSync) {
	// 	return echo.NewHTTPError(http.StatusForbidden, api.FeatureVCSSheetSync.AccessErrorMessage())
	// }
	currentPrincipalID := ctx.Value(common.PrincipalIDContextKey).(int)

	projectResourceID, err := getProjectID(request.Parent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	project, err := s.store.GetProjectV2(ctx, &store.FindProjectMessage{ResourceID: &projectResourceID})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get project with resource id %q, err: %s", projectResourceID, err.Error()))
	}
	if project == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q not found", projectResourceID))
	}
	if project.Deleted {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("project with resource id %q had deleted", projectResourceID))
	}
	if project.Workflow != api.VCSWorkflow {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("project with resource id %q is not a VCS enabled project", projectResourceID))
	}

	repo, err := s.store.GetRepository(ctx, &api.RepositoryFind{ProjectID: &project.UID})
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to find repository for sync sheet: %d", project.UID))
	}
	if repo == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("repository not found for sync sheet: %d", project.UID))
	}

	vcs, err := s.store.GetVCSByID(ctx, repo.VCSID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to find VCS for sync sheet, VCSID: %d", repo.VCSID))
	}
	if vcs == nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("VCS not found for sync sheet: %d", repo.VCSID))
	}

	basePath := filepath.Dir(repo.SheetPathTemplate)
	// TODO(Steven): The repo.branchFilter could be `test/*` which cannot be the ref value.
	fileList, err := vcsPlugin.Get(vcs.Type, vcsPlugin.ProviderConfig{}).FetchRepositoryFileList(ctx,
		common.OauthContext{
			ClientID:     vcs.ApplicationID,
			ClientSecret: vcs.Secret,
			AccessToken:  repo.AccessToken,
			RefreshToken: repo.RefreshToken,
			Refresher:    utils.RefreshToken(ctx, s.store, repo.WebURL),
		},
		vcs.InstanceURL,
		repo.ExternalID,
		repo.BranchFilter,
		basePath,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to fetch repository file list for sync sheet: %d", project.UID))
	}

	for _, file := range fileList {
		sheetInfo, err := parseSheetInfo(file.Path, repo.SheetPathTemplate)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to parse sheet info from template")
		}
		if sheetInfo.SheetName == "" {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("sheet name cannot be empty from sheet path %s with template %s", file.Path, repo.SheetPathTemplate))
		}

		fileContent, err := vcsPlugin.Get(vcs.Type, vcsPlugin.ProviderConfig{}).ReadFileContent(ctx,
			common.OauthContext{
				ClientID:     vcs.ApplicationID,
				ClientSecret: vcs.Secret,
				AccessToken:  repo.AccessToken,
				RefreshToken: repo.RefreshToken,
				Refresher:    utils.RefreshToken(ctx, s.store, repo.WebURL),
			},
			vcs.InstanceURL,
			repo.ExternalID,
			file.Path,
			repo.BranchFilter,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to fetch file content from VCS, instance URL: %s, repo ID: %s, file path: %s, branch: %s", vcs.InstanceURL, repo.ExternalID, file.Path, repo.BranchFilter))
		}

		fileMeta, err := vcsPlugin.Get(vcs.Type, vcsPlugin.ProviderConfig{}).ReadFileMeta(ctx,
			common.OauthContext{
				ClientID:     vcs.ApplicationID,
				ClientSecret: vcs.Secret,
				AccessToken:  repo.AccessToken,
				RefreshToken: repo.RefreshToken,
				Refresher:    utils.RefreshToken(ctx, s.store, repo.WebURL),
			},
			vcs.InstanceURL,
			repo.ExternalID,
			file.Path,
			repo.BranchFilter,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to fetch file meta from VCS, instance URL: %s, repo ID: %s, file path: %s, branch: %s", vcs.InstanceURL, repo.ExternalID, file.Path, repo.BranchFilter))
		}

		lastCommit, err := vcsPlugin.Get(vcs.Type, vcsPlugin.ProviderConfig{}).FetchCommitByID(ctx,
			common.OauthContext{
				ClientID:     vcs.ApplicationID,
				ClientSecret: vcs.Secret,
				AccessToken:  repo.AccessToken,
				RefreshToken: repo.RefreshToken,
				Refresher:    utils.RefreshToken(ctx, s.store, repo.WebURL),
			},
			vcs.InstanceURL,
			repo.ExternalID,
			fileMeta.LastCommitID,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to fetch commit data from VCS, instance URL: %s, repo ID: %s, commit ID: %s", vcs.InstanceURL, repo.ExternalID, fileMeta.LastCommitID))
		}

		sheetVCSPayload := &api.SheetVCSPayload{
			FileName:     fileMeta.Name,
			FilePath:     fileMeta.Path,
			Size:         fileMeta.Size,
			Author:       lastCommit.AuthorName,
			LastCommitID: lastCommit.ID,
			LastSyncTs:   time.Now().Unix(),
		}
		payload, err := json.Marshal(sheetVCSPayload)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to marshal sheetVCSPayload")
		}

		var databaseID *int
		// In non-tenant mode, we can set a databaseId for sheet with ENV_ID and DB_NAME,
		// and ENV_ID and DB_NAME is either both present or neither present.
		if project.TenantMode != api.TenantModeDisabled {
			if sheetInfo.EnvironmentID != "" && sheetInfo.DatabaseName != "" {
				databases, err := s.store.ListDatabases(ctx, &store.FindDatabaseMessage{ProjectID: &project.ResourceID, DatabaseName: &sheetInfo.DatabaseName})
				if err != nil {
					return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to find database list with name: %s, project ID: %d", sheetInfo.DatabaseName, project.UID))
				}
				for _, database := range databases {
					database := database // create a new var "database".
					if database.EnvironmentID == sheetInfo.EnvironmentID {
						databaseID = &database.UID
						break
					}
				}
			}
		}

		var sheetSource api.SheetSource
		switch vcs.Type {
		case vcsPlugin.GitLab:
			sheetSource = api.SheetFromGitLab
		case vcsPlugin.GitHub:
			sheetSource = api.SheetFromGitHub
		case vcsPlugin.Bitbucket:
			sheetSource = api.SheetFromBitbucket
		}
		vscSheetType := api.SheetForSQL
		sheetFind := &api.SheetFind{
			Name:      &sheetInfo.SheetName,
			ProjectID: &project.UID,
			Source:    &sheetSource,
			Type:      &vscSheetType,
		}
		sheet, err := s.store.GetSheet(ctx, sheetFind, currentPrincipalID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to find sheet with name: %s, project ID: %d", sheetInfo.SheetName, project.UID))
		}

		if sheet == nil {
			sheetCreate := api.SheetCreate{
				ProjectID:  project.UID,
				CreatorID:  currentPrincipalID,
				Name:       sheetInfo.SheetName,
				Statement:  fileContent,
				Visibility: api.ProjectSheet,
				Source:     sheetSource,
				Type:       api.SheetForSQL,
				Payload:    string(payload),
			}
			if databaseID != nil {
				sheetCreate.DatabaseID = databaseID
			}

			if _, err := s.store.CreateSheet(ctx, &sheetCreate); err != nil {
				return nil, status.Errorf(codes.Internal, "Failed to create sheet from VCS")
			}
		} else {
			payloadString := string(payload)
			sheetPatch := api.SheetPatch{
				ID:        sheet.ID,
				UpdaterID: currentPrincipalID,
				Statement: &fileContent,
				Payload:   &payloadString,
			}
			if databaseID != nil {
				sheetPatch.DatabaseID = databaseID
			}

			if _, err := s.store.PatchSheet(ctx, &sheetPatch); err != nil {
				return nil, status.Errorf(codes.Internal, "Failed to patch sheet from VCS")
			}
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *SheetService) convertToAPISheetMessage(ctx context.Context, sheet *store.SheetMessage) (*v1pb.Sheet, error) {
	databaseParent := ""
	if sheet.DatabaseID != nil {
		database, err := s.store.GetDatabaseV2(ctx, &store.FindDatabaseMessage{
			UID: sheet.DatabaseID,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get database: %v", err))
		}
		if database == nil {
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("database with id %d not found", *sheet.DatabaseID))
		}
		databaseParent = fmt.Sprintf("%s%s/%s%d", instanceNamePrefix, database.InstanceID, databaseIDPrefix, database.UID)
	}

	visibility := v1pb.Sheet_VISIBILITY_UNSPECIFIED
	switch sheet.Visibility {
	case api.PublicSheet:
		visibility = v1pb.Sheet_VISIBILITY_PUBLIC
	case api.ProjectSheet:
		visibility = v1pb.Sheet_VISIBILITY_PROJECT
	case api.PrivateSheet:
		visibility = v1pb.Sheet_VISIBILITY_PRIVATE
	}

	source := v1pb.Sheet_SOURCE_UNSPECIFIED
	switch sheet.Source {
	case api.SheetFromBytebase:
		source = v1pb.Sheet_SOURCE_BYTEBASE
	case api.SheetFromBytebaseArtifact:
		source = v1pb.Sheet_SOURCE_BYTEBASE_ARTIFACT
	case api.SheetFromGitLab:
		source = v1pb.Sheet_SOURCE_GITLAB
	case api.SheetFromGitHub:
		source = v1pb.Sheet_SOURCE_GITHUB
	case api.SheetFromBitbucket:
		source = v1pb.Sheet_SOURCE_BITBUCKET
	}

	tp := v1pb.Sheet_TYPE_UNSPECIFIED
	switch sheet.Type {
	case api.SheetForSQL:
		tp = v1pb.Sheet_TYPE_SQL
	default:
	}

	creator, err := s.store.GetUserByID(ctx, sheet.CreatorID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get creator: %v", err))
	}

	project, err := s.store.GetProjectByID(ctx, sheet.ProjectUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get project: %v", err))
	}

	return &v1pb.Sheet{
		Name:        fmt.Sprintf("%s%s/%s%d", projectNamePrefix, project.ResourceID, sheetIDPrefix, sheet.UID),
		Database:    databaseParent,
		Title:       sheet.Name,
		Creator:     fmt.Sprintf("users/%s", creator.Email),
		CreateTime:  timestamppb.New(sheet.CreatedTime),
		UpdateTime:  timestamppb.New(sheet.UpdatedTime),
		Content:     []byte(sheet.Statement),
		ContentSize: sheet.Size,
		Visibility:  visibility,
		Source:      source,
		Type:        tp,
		Starred:     sheet.Starred,
		Payload:     sheet.Payload,
	}, nil
}

func convertToStoreSheetMessage(projectUID int, databaseUID *int, creatorID int, sheet *v1pb.Sheet) (*store.SheetMessage, error) {
	visibility, err := convertToLegacyAPISheetVisibility(sheet.Visibility)
	if err != nil {
		return nil, err
	}
	var source api.SheetSource
	switch sheet.Source {
	case v1pb.Sheet_SOURCE_UNSPECIFIED:
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid source %q", sheet.Source))
	case v1pb.Sheet_SOURCE_BYTEBASE:
		source = api.SheetFromBytebase
	case v1pb.Sheet_SOURCE_BYTEBASE_ARTIFACT:
		source = api.SheetFromBytebaseArtifact
	case v1pb.Sheet_SOURCE_GITLAB:
		source = api.SheetFromGitLab
	case v1pb.Sheet_SOURCE_GITHUB:
		source = api.SheetFromGitHub
	case v1pb.Sheet_SOURCE_BITBUCKET:
		source = api.SheetFromBitbucket
	default:
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid source %q", sheet.Source))
	}
	var tp api.SheetType
	switch sheet.Type {
	case v1pb.Sheet_TYPE_UNSPECIFIED:
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid type %q", sheet.Type))
	case v1pb.Sheet_TYPE_SQL:
		tp = api.SheetForSQL
	default:
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid type %q", sheet.Type))
	}

	return &store.SheetMessage{
		ProjectUID: projectUID,
		DatabaseID: databaseUID,
		CreatorID:  creatorID,
		Name:       sheet.Title,
		Statement:  string(sheet.Content),
		Visibility: visibility,
		Source:     source,
		Type:       tp,
		Payload:    sheet.Payload,
	}, nil
}

func convertToLegacyAPISheetVisibility(visibility v1pb.Sheet_Visibility) (api.SheetVisibility, error) {
	switch visibility {
	case v1pb.Sheet_VISIBILITY_UNSPECIFIED:
		return api.SheetVisibility(""), status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid visibility %q", visibility))
	case v1pb.Sheet_VISIBILITY_PUBLIC:
		return api.PublicSheet, nil
	case v1pb.Sheet_VISIBILITY_PROJECT:
		return api.ProjectSheet, nil
	case v1pb.Sheet_VISIBILITY_PRIVATE:
		return api.PrivateSheet, nil
	default:
		return api.SheetVisibility(""), status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid visibility %q", visibility))
	}
}

// SheetInfo represents the sheet related information from sheetPathTemplate.
type SheetInfo struct {
	EnvironmentID string
	DatabaseName  string
	SheetName     string
}

// parseSheetInfo matches sheetPath against sheetPathTemplate. If sheetPath matches, then it will derive SheetInfo from the sheetPath.
// Both sheetPath and sheetPathTemplate are the full file path(including the base directory) of the repository.
func parseSheetInfo(sheetPath string, sheetPathTemplate string) (*SheetInfo, error) {
	placeholderList := []string{
		"ENV_ID",
		"DB_NAME",
		"NAME",
	}
	sheetPathRegex := sheetPathTemplate
	for _, placeholder := range placeholderList {
		sheetPathRegex = strings.ReplaceAll(sheetPathRegex, fmt.Sprintf("{{%s}}", placeholder), fmt.Sprintf("(?P<%s>[a-zA-Z0-9\\+\\-\\=\\_\\#\\!\\$\\. ]+)", placeholder))
	}
	sheetRegex, err := regexp.Compile(fmt.Sprintf("^%s$", sheetPathRegex))
	if err != nil {
		return nil, errors.Wrapf(err, "invalid sheet path template %q", sheetPathTemplate)
	}
	if !sheetRegex.MatchString(sheetPath) {
		return nil, errors.Errorf("sheet path %q does not match sheet path template %q", sheetPath, sheetPathTemplate)
	}

	matchList := sheetRegex.FindStringSubmatch(sheetPath)
	sheetInfo := &SheetInfo{}
	for _, placeholder := range placeholderList {
		index := sheetRegex.SubexpIndex(placeholder)
		if index >= 0 {
			switch placeholder {
			case "ENV_ID":
				sheetInfo.EnvironmentID = matchList[index]
			case "DB_NAME":
				sheetInfo.DatabaseName = matchList[index]
			case "NAME":
				sheetInfo.SheetName = matchList[index]
			}
		}
	}

	return sheetInfo, nil
}