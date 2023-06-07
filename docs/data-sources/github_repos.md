---
page_title: "autocloud_github_repos Data Source - AutoCloud"
subcategory: ""
description: |-
  A data resource to describe and work with the Github repositories connected to this AutoCloud organization.
---

# Data Source: autocloud_github_repos

A data resource to describe and work with the Github repositories connected to this AutoCloud organization, for connecting the code delivery output of an AutoCloud infrastructure as code blueprint.

For more details on connecting AutoCloud to your Github repositories, please refer to the [source control integration documentation](https://docs.autocloud.io/integration-with-source-control-github).


## Example Usage

```terraform
data "autocloud_github_repos" "repos" {}

####
# Local variables
locals {
  # Destination repos where generated code will be submitted
  dest_repos = [
    for repo in data.autocloud_github_repos.repos.data[*].url : repo
    if length(regexall("/infrastructure-live-demo", repo)) > 0
  ]
}

####
# Create Blueprint
#
# Create generator blueprint that contains all the elements
resource "autocloud_blueprint" "this" {
  name = "KMS Encrypted S3 Bucket"

  .
  .
  .

  ###
  # Destination repository git configuraiton
  #
  git_config {
    destination_branch = "main"

    git_url_options = local.dest_repos
    git_url_default = length(local.dest_repos) != 0 ? local.dest_repos[0] : "" # Choose the first in the list by default
  }
}
```

## Arguments

None


## Attributes

- [`data`](#nestedatt--data) - (List of Object) The Github repository information.
- `id` - (String) The ID of this resource.

<a id="nestedatt--data"></a>
### data

- `description` - (String) The Github repository description.
- `id` - (Number) The Github repository ID.
- `name` - (String) The Github repository name.
- `url` - (String) The Github repository URL.
