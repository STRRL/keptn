package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

type ExpandedProjectOld struct {

	// Creation date of the project
	CreationDate string `json:"creationDate,omitempty"`

	// Git remote URI
	GitRemoteURI string `json:"gitRemoteURI,omitempty"`

	// Git User
	GitUser string `json:"gitUser,omitempty"`

	// last event context
	LastEventContext *apimodels.EventContextInfo `json:"lastEventContext,omitempty"`

	// Project name
	ProjectName string `json:"projectName,omitempty"`

	// Shipyard file content
	Shipyard string `json:"shipyard,omitempty"`

	// Version of the shipyard file
	ShipyardVersion string `json:"shipyardVersion,omitempty"`

	// git proxy URL
	GitProxyURL string `json:"gitProxyUrl,omitempty"`

	// git proxy scheme
	GitProxyScheme string `json:"gitProxyScheme,omitempty"`

	// git proxy user
	GitProxyUser string `json:"gitProxyUser,omitempty"`

	// insecure skip tls
	InsecureSkipTLS bool `json:"insecureSkipTLS"`

	// stages
	Stages []*apimodels.ExpandedStage `json:"stages"`
}

type MongoDBProjectCredentialsRepo struct {
	ProjectRepo *MongoDBProjectsRepo
}

func NewMongoDBProjectCredentialsRepo(dbConnection *MongoDBConnection) *MongoDBProjectCredentialsRepo {
	projectsRepo := NewMongoDBProjectsRepo(dbConnection)
	return &MongoDBProjectCredentialsRepo{
		ProjectRepo: projectsRepo,
	}
}

func (m *MongoDBProjectCredentialsRepo) GetOldCredentialsProjects() ([]*ExpandedProjectOld, error) {
	result := []*ExpandedProjectOld{}
	err := m.ProjectRepo.DBConnection.EnsureDBConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projectCollection := m.ProjectRepo.getProjectsCollection()
	cursor, err := projectCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error retrieving projects from mongoDB: " + err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		projectResult := &ExpandedProjectOld{}
		err := cursor.Decode(projectResult)
		if err != nil {
			fmt.Println("Could not cast to *models.Project")
		}
		result = append(result, projectResult)
	}

	return result, nil
}

func TransformGitCredentials(project *ExpandedProjectOld) *apimodels.ExpandedProject {
	//if project has no credentials, or has credentials in the newest format
	if project.GitRemoteURI == "" {
		return nil
	}

	newProject := apimodels.ExpandedProject{}
	newProject.CreationDate = project.CreationDate
	newProject.LastEventContext = project.LastEventContext
	newProject.ProjectName = project.ProjectName
	newProject.Shipyard = project.Shipyard
	newProject.ShipyardVersion = project.ShipyardVersion
	newProject.Stages = project.Stages

	//project has credentials in old format
	credentials := apimodels.GitAuthCredentialsSecure{
		RemoteURL: project.GitRemoteURI,
		User:      project.GitUser,
	}
	newProject.GitCredentials = &credentials

	//if project is using ssh auth, no other parameters are stored
	if strings.HasPrefix(project.GitRemoteURI, "ssh://") {
		return &newProject
	}

	//project is using https auth, InsecureSkipTLS needs to be set
	httpCredentials := apimodels.HttpsGitAuthSecure{
		InsecureSkipTLS: project.InsecureSkipTLS,
	}
	newProject.GitCredentials.HttpsAuth = &httpCredentials

	//project is not using proxy, no additional parameters need to be stored
	if project.GitProxyURL == "" {
		return &newProject
	}

	//project is using proxy
	proxyCredentials := apimodels.ProxyGitAuthSecure{
		Scheme: project.GitProxyScheme,
		URL:    project.GitProxyURL,
	}
	newProject.GitCredentials.HttpsAuth.Proxy = &proxyCredentials

	//project is using proxy with a user
	if project.GitProxyUser != "" {
		newProject.GitCredentials.HttpsAuth.Proxy.User = project.GitProxyUser
	}

	return &newProject
}

func (m *MongoDBProjectCredentialsRepo) UpdateProject(project *ExpandedProjectOld) error {
	newProject := TransformGitCredentials(project)
	if newProject == nil {
		return nil
	}

	return m.ProjectRepo.UpdateProject(newProject)
}
