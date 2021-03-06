// Copyright 2018 Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/goharbor/harbor/src/core/api"
	"github.com/goharbor/harbor/src/core/config"
	"github.com/goharbor/harbor/src/core/controllers"
	"github.com/goharbor/harbor/src/core/service/notifications/admin"
	"github.com/goharbor/harbor/src/core/service/notifications/clair"
	"github.com/goharbor/harbor/src/core/service/notifications/jobs"
	"github.com/goharbor/harbor/src/core/service/notifications/registry"
	"github.com/goharbor/harbor/src/core/service/token"

	"github.com/astaxie/beego"
)

func initRouters() {

	// standalone
	if !config.WithAdmiral() {
		// Controller API:
		beego.Router("/c/login", &controllers.CommonController{}, "post:Login")
		beego.Router("/c/log_out", &controllers.CommonController{}, "get:LogOut")
		beego.Router("/c/reset", &controllers.CommonController{}, "post:ResetPassword")
		beego.Router("/c/userExists", &controllers.CommonController{}, "post:UserExists")
		beego.Router("/c/sendEmail", &controllers.CommonController{}, "get:SendResetEmail")

		// API:
		beego.Router("/api/projects/:pid([0-9]+)/members/?:pmid([0-9]+)", &api.ProjectMemberAPI{})
		beego.Router("/api/projects/", &api.ProjectAPI{}, "head:Head")
		beego.Router("/api/projects/:id([0-9]+)", &api.ProjectAPI{})

		beego.Router("/api/users/:id", &api.UserAPI{}, "get:Get;delete:Delete;put:Put")
		beego.Router("/api/users", &api.UserAPI{}, "get:List;post:Post")
		beego.Router("/api/users/:id([0-9]+)/password", &api.UserAPI{}, "put:ChangePassword")
		beego.Router("/api/users/:id/sysadmin", &api.UserAPI{}, "put:ToggleUserAdminRole")
		beego.Router("/api/usergroups/?:ugid([0-9]+)", &api.UserGroupAPI{})
		beego.Router("/api/ldap/ping", &api.LdapAPI{}, "post:Ping")
		beego.Router("/api/ldap/users/search", &api.LdapAPI{}, "get:Search")
		beego.Router("/api/ldap/groups/search", &api.LdapAPI{}, "get:SearchGroup")
		beego.Router("/api/ldap/users/import", &api.LdapAPI{}, "post:ImportUser")
		beego.Router("/api/email/ping", &api.EmailAPI{}, "post:Ping")
	}

	// API
	beego.Router("/api/health", &api.HealthAPI{}, "get:CheckHealth")
	beego.Router("/api/ping", &api.SystemInfoAPI{}, "get:Ping")
	beego.Router("/api/search", &api.SearchAPI{})
	beego.Router("/api/projects/", &api.ProjectAPI{}, "get:List;post:Post")
	beego.Router("/api/projects/:id([0-9]+)/logs", &api.ProjectAPI{}, "get:Logs")
	beego.Router("/api/projects/:id([0-9]+)/_deletable", &api.ProjectAPI{}, "get:Deletable")
	beego.Router("/api/projects/:id([0-9]+)/metadatas/?:name", &api.MetadataAPI{}, "get:Get")
	beego.Router("/api/projects/:id([0-9]+)/metadatas/", &api.MetadataAPI{}, "post:Post")
	beego.Router("/api/projects/:id([0-9]+)/metadatas/:name", &api.MetadataAPI{}, "put:Put;delete:Delete")
	beego.Router("/api/repositories", &api.RepositoryAPI{}, "get:Get")
	beego.Router("/api/repositories/scanAll", &api.RepositoryAPI{}, "post:ScanAll")
	beego.Router("/api/repositories/*", &api.RepositoryAPI{}, "delete:Delete;put:Put")
	beego.Router("/api/repositories/*/labels", &api.RepositoryLabelAPI{}, "get:GetOfRepository;post:AddToRepository")
	beego.Router("/api/repositories/*/labels/:id([0-9]+)", &api.RepositoryLabelAPI{}, "delete:RemoveFromRepository")
	beego.Router("/api/repositories/*/tags/:tag", &api.RepositoryAPI{}, "delete:Delete;get:GetTag")
	beego.Router("/api/repositories/*/tags/:tag/labels", &api.RepositoryLabelAPI{}, "get:GetOfImage;post:AddToImage")
	beego.Router("/api/repositories/*/tags/:tag/labels/:id([0-9]+)", &api.RepositoryLabelAPI{}, "delete:RemoveFromImage")
	beego.Router("/api/repositories/*/tags", &api.RepositoryAPI{}, "get:GetTags;post:Retag")
	beego.Router("/api/repositories/*/tags/:tag/scan", &api.RepositoryAPI{}, "post:ScanImage")
	beego.Router("/api/repositories/*/tags/:tag/vulnerability/details", &api.RepositoryAPI{}, "Get:VulnerabilityDetails")
	beego.Router("/api/repositories/*/tags/:tag/manifest", &api.RepositoryAPI{}, "get:GetManifests")
	beego.Router("/api/repositories/*/signatures", &api.RepositoryAPI{}, "get:GetSignatures")
	beego.Router("/api/repositories/top", &api.RepositoryAPI{}, "get:GetTopRepos")
	beego.Router("/api/jobs/replication/", &api.RepJobAPI{}, "get:List;put:StopJobs")
	beego.Router("/api/jobs/replication/:id([0-9]+)", &api.RepJobAPI{})
	beego.Router("/api/jobs/replication/:id([0-9]+)/log", &api.RepJobAPI{}, "get:GetLog")
	beego.Router("/api/jobs/scan/:id([0-9]+)/log", &api.ScanJobAPI{}, "get:GetLog")

	beego.Router("/api/system/gc", &api.GCAPI{}, "get:List")
	beego.Router("/api/system/gc/:id", &api.GCAPI{}, "get:GetGC")
	beego.Router("/api/system/gc/:id([0-9]+)/log", &api.GCAPI{}, "get:GetLog")
	beego.Router("/api/system/gc/schedule", &api.GCAPI{}, "get:Get;put:Put;post:Post")

	beego.Router("/api/policies/replication/:id([0-9]+)", &api.RepPolicyAPI{})
	beego.Router("/api/policies/replication", &api.RepPolicyAPI{}, "get:List")
	beego.Router("/api/policies/replication", &api.RepPolicyAPI{}, "post:Post")
	beego.Router("/api/targets/", &api.TargetAPI{}, "get:List")
	beego.Router("/api/targets/", &api.TargetAPI{}, "post:Post")
	beego.Router("/api/targets/:id([0-9]+)", &api.TargetAPI{})
	beego.Router("/api/targets/:id([0-9]+)/policies/", &api.TargetAPI{}, "get:ListPolicies")
	beego.Router("/api/targets/ping", &api.TargetAPI{}, "post:Ping")
	beego.Router("/api/logs", &api.LogAPI{})
	beego.Router("/api/configs", &api.ConfigAPI{}, "get:GetInternalConfig")
	beego.Router("/api/configurations", &api.ConfigAPI{})
	beego.Router("/api/configurations/reset", &api.ConfigAPI{}, "post:Reset")
	beego.Router("/api/statistics", &api.StatisticAPI{})
	beego.Router("/api/replications", &api.ReplicationAPI{})
	beego.Router("/api/labels", &api.LabelAPI{}, "post:Post;get:List")
	beego.Router("/api/labels/:id([0-9]+)", &api.LabelAPI{}, "get:Get;put:Put;delete:Delete")
	beego.Router("/api/labels/:id([0-9]+)/resources", &api.LabelAPI{}, "get:ListResources")

	beego.Router("/api/systeminfo", &api.SystemInfoAPI{}, "get:GetGeneralInfo")
	beego.Router("/api/systeminfo/volumes", &api.SystemInfoAPI{}, "get:GetVolumeInfo")
	beego.Router("/api/systeminfo/getcert", &api.SystemInfoAPI{}, "get:GetCert")

	beego.Router("/api/internal/syncregistry", &api.InternalAPI{}, "post:SyncRegistry")
	beego.Router("/api/internal/renameadmin", &api.InternalAPI{}, "post:RenameAdmin")
	beego.Router("/api/internal/configurations", &api.ConfigAPI{}, "get:GetInternalConfig")

	// external service that hosted on harbor process:
	beego.Router("/service/notifications", &registry.NotificationHandler{})
	beego.Router("/service/notifications/clair", &clair.Handler{}, "post:Handle")
	beego.Router("/service/notifications/jobs/scan/:id([0-9]+)", &jobs.Handler{}, "post:HandleScan")
	beego.Router("/service/notifications/jobs/replication/:id([0-9]+)", &jobs.Handler{}, "post:HandleReplication")
	beego.Router("/service/notifications/jobs/adminjob/:id([0-9]+)", &admin.Handler{}, "post:HandleAdminJob")
	beego.Router("/service/token", &token.Handler{})

	beego.Router("/v2/*", &controllers.RegistryProxy{}, "*:Handle")

	// APIs for chart repository
	if config.WithChartMuseum() {
		// Charts are controlled under projects
		chartRepositoryAPIType := &api.ChartRepositoryAPI{}
		beego.Router("/api/chartrepo/health", chartRepositoryAPIType, "get:GetHealthStatus")
		beego.Router("/api/chartrepo/:repo/charts", chartRepositoryAPIType, "get:ListCharts")
		beego.Router("/api/chartrepo/:repo/charts/:name", chartRepositoryAPIType, "get:ListChartVersions")
		beego.Router("/api/chartrepo/:repo/charts/:name", chartRepositoryAPIType, "delete:DeleteChart")
		beego.Router("/api/chartrepo/:repo/charts/:name/:version", chartRepositoryAPIType, "get:GetChartVersion")
		beego.Router("/api/chartrepo/:repo/charts/:name/:version", chartRepositoryAPIType, "delete:DeleteChartVersion")
		beego.Router("/api/chartrepo/:repo/charts", chartRepositoryAPIType, "post:UploadChartVersion")
		beego.Router("/api/chartrepo/:repo/prov", chartRepositoryAPIType, "post:UploadChartProvFile")
		beego.Router("/api/chartrepo/charts", chartRepositoryAPIType, "post:UploadChartVersion")

		// Repository services
		beego.Router("/chartrepo/:repo/index.yaml", chartRepositoryAPIType, "get:GetIndexByRepo")
		beego.Router("/chartrepo/index.yaml", chartRepositoryAPIType, "get:GetIndex")
		beego.Router("/chartrepo/:repo/charts/:filename", chartRepositoryAPIType, "get:DownloadChart")

		// Labels for chart
		chartLabelAPIType := &api.ChartLabelAPI{}
		beego.Router("/api/chartrepo/:repo/charts/:name/:version/labels", chartLabelAPIType, "get:GetLabels;post:MarkLabel")
		beego.Router("/api/chartrepo/:repo/charts/:name/:version/labels/:id([0-9]+)", chartLabelAPIType, "delete:RemoveLabel")
	}

	// Error pages
	beego.ErrorController(&controllers.ErrorController{})

}
