package restapi

import (
	"testing"
)

func parseRoleQueryResponseString(yaml string) (RoleQueryResponse, error) {
	return ParseRoleQueryResponse([]byte(yaml))
}

func TestParseRoleQueryResponse(t *testing.T) {
	roleQueryResponse, err := parseRoleQueryResponseString(`{
		"count": 1,
		"cur_page": 1,
		"num_pages": 1,
		"next_link": null,
		"previous_link": null,
		"next": null,
		"previous": null,
		"results": [
			{
				"url": "/api/v1/roles/10823/",
				"related": {
					"imports": "/api/v1/roles/10823/imports/",
					"versions": "/api/v1/roles/10823/versions/",
					"dependencies": "/api/v1/roles/10823/dependencies/",
					"notifications": "/api/v1/roles/10823/notifications/"
				},
				"summary_fields": {
					"platforms": [
						{
							"release": "trusty",
							"name": "Ubuntu"
						},
						{
							"release": "wily",
							"name": "Ubuntu"
						},
						{
							"release": "xenial",
							"name": "Ubuntu"
						}
					],
					"versions": [
						{
							"release_date": "2016-09-05T20:42:09Z",
							"name": "1.1.3",
							"id": 18475
						},
						{
							"release_date": "2016-08-25T02:26:01Z",
							"name": "1.1.2",
							"id": 18144
						},
						{
							"release_date": "2016-08-23T15:09:43Z",
							"name": "1.1.1",
							"id": 18034
						},
						{
							"release_date": "2016-08-16T19:54:05Z",
							"name": "1.1.0",
							"id": 17837
						},
						{
							"release_date": "2016-07-08T13:52:32Z",
							"name": "1.0.0",
							"id": 16412
						}
					],
					"dependencies": [],
					"tags": [
						{
							"name": "apt"
						}
					]
				},
				"id": 10823,
				"created": "2016-07-08T09:57:51.931Z",
				"modified": "2016-12-06T14:57:35.389Z",
				"name": "apt",
				"role_type": "ANS",
				"namespace": "gantsign",
				"is_valid": true,
				"github_user": "gantsign",
				"github_repo": "ansible-role-apt",
				"github_branch": "",
				"min_ansible_version": "1.9",
				"issue_tracker_url": "https://github.com/gantsign/ansible-role-apt/issues",
				"license": "MIT",
				"company": "GantSign Ltd.",
				"description": "Role for configuring the APT package manager.",
				"readme": "Ansible Role: APT\n=================\n\n[![Build Status](https://travis-ci.org/gantsign/ansible-role-apt.svg?branch=master)](https://travis-ci.org/gantsign/ansible-role-apt)\n[![Ansible Galaxy](https://img.shields.io/badge/ansible--galaxy-gantsign.apt-blue.svg)](https://galaxy.ansible.com/gantsign/apt)\n[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/gantsign/ansible-role-apt/master/LICENSE)\n\nRole to configure the APT package manager. Currently limited to controlling the\nproperties that affect the cleaning of the DEB files (typically by the APT cron\njob). The DEB files are removed to save on disk space but if you're using\nVagrant (with the vagrant-cachier plugin) you may want to keep the DEB files to\nspeed up VM rebuilds.\n\nRequirements\n------------\n\n* Ubuntu\n\nRole Variables\n--------------\n\nThe following variables will change the behavior of this role (default values\nare shown below):\n\n~~~yaml\n# The filename for the apt config\napt_config_filename: '80-vagrant'\n\n# Whether the cache of DEB files should be preserved or cleaned\napt_preserve_cache: no\n\n# Max age (in days) of DEB files to keep when cleaning cache\napt_archives_maxage: null\n\n# Min age (in days) of DEB files to keep when cleaning cache\napt_archives_minage: null\n\n# Max size (in MB) of DEB files to keep when cleaning cache\napt_archives_maxsize: null\n~~~\n\nExample Playbook\n----------------\n\n~~~yaml\n- hosts: servers\n  roles:\n     - { role: gantsign.apt, apt_preserve_cache: yes }\n~~~\n\nMore Roles From GantSign\n------------------------\n\nYou can find more roles from GantSign on\n[Ansible Galaxy](https://galaxy.ansible.com/gantsign).\n\nDevelopment & Testing\n---------------------\n\nThis project uses [Molecule](http://molecule.readthedocs.io/) to aid in the\ndevelopment and testing; the role is unit tested using\n[Testinfra](http://testinfra.readthedocs.io/) and\n[pytest](http://docs.pytest.org/).\n\nTo develop or test you'll need to have installed the following:\n\n* Linux (e.g. [Ubuntu](http://www.ubuntu.com/))\n* [Docker](https://www.docker.com/)\n* [Python](https://www.python.org/) (including python-pip)\n* [Ansible](https://www.ansible.com/)\n* [Molecule](http://molecule.readthedocs.io/)\n\nTo run the role (i.e. the tests/test.yml playbook), and test the results\n(tests/test_role.py), execute the following command from the project root\n(i.e. the directory with molecule.yml in it):\n\n~~~bash\nmolecule test\n~~~\n\nLicense\n-------\n\nMIT\n\nAuthor Information\n------------------\n\nJohn Freeman\n\nGantSign Ltd.\nCompany No. 06109112 (registered in England)\n",
				"readme_html": "<h1>Ansible Role: APT<\/h1>\n<p><a href=\"https://travis-ci.org/gantsign/ansible-role-apt\"><img alt=\"Build Status\" src=\"https://travis-ci.org/gantsign/ansible-role-apt.svg?branch=master\" /><\/a>\n<a href=\"https://galaxy.ansible.com/gantsign/apt\"><img alt=\"Ansible Galaxy\" src=\"https://img.shields.io/badge/ansible--galaxy-gantsign.apt-blue.svg\" /><\/a>\n<a href=\"https://raw.githubusercontent.com/gantsign/ansible-role-apt/master/LICENSE\"><img alt=\"License\" src=\"https://img.shields.io/badge/license-MIT-blue.svg\" /><\/a><\/p>\n<p>Role to configure the APT package manager. Currently limited to controlling the\nproperties that affect the cleaning of the DEB files (typically by the APT cron\njob). The DEB files are removed to save on disk space but if you're using\nVagrant (with the vagrant-cachier plugin) you may want to keep the DEB files to\nspeed up VM rebuilds.<\/p>\n<h2>Requirements<\/h2>\n<ul>\n<li>Ubuntu<\/li>\n<\/ul>\n<h2>Role Variables<\/h2>\n<p>The following variables will change the behavior of this role (default values\nare shown below):<\/p>\n<pre><code class=\"yaml\"># The filename for the apt config\napt_config_filename: '80-vagrant'\n\n# Whether the cache of DEB files should be preserved or cleaned\napt_preserve_cache: no\n\n# Max age (in days) of DEB files to keep when cleaning cache\napt_archives_maxage: null\n\n# Min age (in days) of DEB files to keep when cleaning cache\napt_archives_minage: null\n\n# Max size (in MB) of DEB files to keep when cleaning cache\napt_archives_maxsize: null\n<\/code><\/pre>\n\n<h2>Example Playbook<\/h2>\n<pre><code class=\"yaml\">- hosts: servers\n  roles:\n     - { role: gantsign.apt, apt_preserve_cache: yes }\n<\/code><\/pre>\n\n<h2>More Roles From GantSign<\/h2>\n<p>You can find more roles from GantSign on\n<a href=\"https://galaxy.ansible.com/gantsign\">Ansible Galaxy<\/a>.<\/p>\n<h2>Development & Testing<\/h2>\n<p>This project uses <a href=\"http://molecule.readthedocs.io/\">Molecule<\/a> to aid in the\ndevelopment and testing; the role is unit tested using\n<a href=\"http://testinfra.readthedocs.io/\">Testinfra<\/a> and\n<a href=\"http://docs.pytest.org/\">pytest<\/a>.<\/p>\n<p>To develop or test you'll need to have installed the following:<\/p>\n<ul>\n<li>Linux (e.g. <a href=\"http://www.ubuntu.com/\">Ubuntu<\/a>)<\/li>\n<li><a href=\"https://www.docker.com/\">Docker<\/a><\/li>\n<li><a href=\"https://www.python.org/\">Python<\/a> (including python-pip)<\/li>\n<li><a href=\"https://www.ansible.com/\">Ansible<\/a><\/li>\n<li><a href=\"http://molecule.readthedocs.io/\">Molecule<\/a><\/li>\n<\/ul>\n<p>To run the role (i.e. the <code>tests/test.yml<\/code> playbook), and test the results\n(<code>tests/test_role.py<\/code>), execute the following command from the project root\n(i.e. the directory with <code>molecule.yml<\/code> in it):<\/p>\n<pre><code class=\"bash\">molecule test\n<\/code><\/pre>\n\n<h2>License<\/h2>\n<p>MIT<\/p>\n<h2>Author Information<\/h2>\n<p>John Freeman<\/p>\n<p>GantSign Ltd.\nCompany No. 06109112 (registered in England)<\/p>",
				"travis_status_url": "https://travis-ci.org/gantsign/ansible-role-apt.svg?branch=master",
				"stargazers_count": 1,
				"watchers_count": 1,
				"forks_count": 0,
				"open_issues_count": 0,
				"commit": "e84f17b6d17c2c76fe8d212ffc799fe36f1dff02",
				"commit_message": "Updated Molecule to 1.11.1 (#28)\n\nKeeping up with the latest changes.",
				"commit_url": "https://github.com/gantsign/ansible-role-apt/commit/e84f17b6d17c2c76fe8d212ffc799fe36f1dff02",
				"download_count": 519,
				"active": true
			}
		]
	}`)
	if err != nil {
		t.Errorf("Error parsing bad YAML: %+v", err)
		return
	}

	expectedResultsLength := 1
	actualResultsLength := len(roleQueryResponse.Results)
	if expectedResultsLength != actualResultsLength {
		t.Errorf("Expected [%d] != [%d]", expectedResultsLength, actualResultsLength)
		return
	}

	roleDetails := roleQueryResponse.Results[0]

	expectedGitHubUser := "gantsign"
	actualGitHubUser := roleDetails.GitHubUser
	if expectedGitHubUser != actualGitHubUser {
		t.Errorf("Expected [%s] != actual [%s]", expectedGitHubUser, actualGitHubUser)
	}

	expectedGitHubRepo := "ansible-role-apt"
	actualGitHubRepo := roleDetails.GitHubRepo
	if expectedGitHubRepo != actualGitHubRepo {
		t.Errorf("Expected [%s] != actual [%s]", expectedGitHubRepo, actualGitHubRepo)
	}

	versions := roleDetails.SummaryFields.Versions

	expectedVersionsLength := 5
	actualVersionsLength := len(versions)
	if expectedVersionsLength != actualVersionsLength {
		t.Errorf("Expected [%d] != [%d]", expectedVersionsLength, actualVersionsLength)
		return
	}

	version := versions[0]

	expectedReleaseDate := "2016-09-05T20:42:09Z"
	actualReleaseDate := version.ReleaseDate
	if expectedReleaseDate != actualReleaseDate {
		t.Errorf("Expected [%s] != actual [%s]", expectedReleaseDate, actualReleaseDate)
	}

	expectedName := "1.1.3"
	actualName := version.Name
	if expectedName != actualName {
		t.Errorf("Expected [%s] != actual [%s]", expectedName, actualName)
	}

	expectedId := 18475
	actualId := version.Id
	if expectedId != actualId {
		t.Errorf("Expected [%d] != [%d]", expectedId, actualId)
	}
}

func TestLatestVersionWithoutPrefix(t *testing.T) {

	roleQueryResponse, err := parseRoleQueryResponseString(`{
		"results": [
			{
				"summary_fields": {
					"versions": [
						{
							"release_date": "2016-09-05T20:42:09Z",
							"name": "1.1.3",
							"id": 18475
						},
						{
							"release_date": "2016-08-25T02:26:01Z",
							"name": "1.1.2",
							"id": 18144
						},
						{
							"release_date": "2016-08-23T15:09:43Z",
							"name": "1.1.1",
							"id": 18034
						},
						{
							"release_date": "2016-08-16T19:54:05Z",
							"name": "1.1.0",
							"id": 17837
						},
						{
							"release_date": "2016-07-08T13:52:32Z",
							"name": "1.0.0",
							"id": 16412
						}
					]
				}
			}
		]
	}`)
	if err != nil {
		t.Errorf("Error parsing bad YAML: %+v", err)
		return
	}

	expectedVersion := "1.1.3"

	roleDetails := roleQueryResponse.Results[0]
	actualVersion := roleDetails.LatestVersion()

	if expectedVersion != actualVersion {
		t.Errorf("Expected [%s] != actual [%s]", expectedVersion, actualVersion)
	}
}

func TestLatestVersionWithPrefix(t *testing.T) {
	roleQueryResponse, err := parseRoleQueryResponseString(`{
		"results": [
			{
				"summary_fields": {
					"versions": [
						{
							"release_date": "2016-09-19T10:07:03Z",
							"name": "v2.0.1",
							"id": 19123
						},
						{
							"release_date": "2016-08-26T13:06:41Z",
							"name": "v2.0.0",
							"id": 18132
						},
						{
							"release_date": "2016-08-23T10:20:49Z",
							"name": "v1.2.1",
							"id": 18050
						},
						{
							"release_date": "2016-08-15T13:15:08Z",
							"name": "v1.2.0",
							"id": 17796
						},
						{
							"release_date": "2016-08-15T12:39:53Z",
							"name": "v1.1.0",
							"id": 17794
						},
						{
							"release_date": "2016-04-20T20:04:01Z",
							"name": "v1.0.3",
							"id": 17795
						},
						{
							"release_date": "2016-04-20T17:46:07Z",
							"name": "v1.0.2",
							"id": 18051
						},
						{
							"release_date": "2016-04-15T13:39:51Z",
							"name": "v1.0.1",
							"id": 13117
						},
						{
							"release_date": "2016-04-15T13:14:31Z",
							"name": "v1.0.0",
							"id": 13116
						},
						{
							"release_date": "2016-08-23T14:20:49Z",
							"name": "1.2.1",
							"id": 18045
						},
						{
							"release_date": "2016-04-20T21:46:07Z",
							"name": "1.0.2",
							"id": 13349
						}
					]
				}
			}
		]
	}`)
	if err != nil {
		t.Errorf("Error parsing bad YAML: %+v", err)
		return
	}

	expectedVersion := "v2.0.1"

	roleDetails := roleQueryResponse.Results[0]
	actualVersion := roleDetails.LatestVersion()

	if expectedVersion != actualVersion {
		t.Errorf("Expected [%s] != actual [%s]", expectedVersion, actualVersion)
	}
}

func TestLatestVersionEmpty(t *testing.T) {

	roleQueryResponse, err := parseRoleQueryResponseString(`{
		"results": [
			{
				"summary_fields": {
					"versions": []
				}
			}
		]
	}`)
	if err != nil {
		t.Errorf("Error parsing bad YAML: %+v", err)
		return
	}

	expectedVersion := "master"

	roleDetails := roleQueryResponse.Results[0]
	actualVersion := roleDetails.LatestVersion()

	if expectedVersion != actualVersion {
		t.Errorf("Expected [%s] != actual [%s]", expectedVersion, actualVersion)
	}
}
