# prm (Pull Request Monitor)
`prm` is a lightweight command-line tool designed to help developers track their pull requests across multiple repositories and Source Code Management (SCM) providers effortlessly.

## Highlights
- **Simple Setup:** Provide minimal information about the SCM provider, and `prm` automatically fetches all your repositories.
- **Fast Performance:** Fetch all your PRs, even if there are hundreds, in just a couple of seconds.
- **Supports JSON & YAML:** Output your data in JSON or YAML for easy integration with other tools.
- **Secure:** All data stays on your local machine, ensuring privacy. You can purge any locally persisted data with a single command.

## What Are Source Code Management Providers?
A Source Code Management (SCM) provider is any software solution that allows you to host Git repositories, such as GitHub, Harness, GitLab, etc.

### Supported SCM Providers
- **GitHub**
- **Harness**

## Requirements
For any SCM provider you want to use, you need to have the following:
- **Host URL:** The base URL of that provider (e.g., `https://github.com` for GitHub or `https://app.harness.io` for Harness). It can be a cloud service or a self-managed instance, as long as it is accessible from your machine.
- **PAT:** The Personal Access Token (PAT) is a secure way to authenticate with your SCM provider. It allows `prm` to access details like user information, pull requests, and reviewers.

## Installation
### Using the installation script
#### Installing the latest version
You can run the following command to download the installation script and install `prm`. It detects the OS and Arch of your laptop and downloads the latest release of the corresponding binary file.
```bash
curl -L https://raw.githubusercontent.com/dhruv1397/prm/main/install.sh | bash
```
<img width="1430" alt="installation" src="https://github.com/user-attachments/assets/c278377b-9d8a-45f9-afdc-24d8b75cfad7">

#### Installing a specific version
To install a specific version, replace target_version with the version number you'd like to install:
```bash
curl -L https://raw.githubusercontent.com/dhruv1397/prm/main/install.sh | bash -s -- target_version
```
### Downloading the binary manually
You can download the binary directly from the release page, make it executable and add it to the PATH.

## Usage
> All flags support shorthands ie -t for --type, -o for --output, -s for --state, -n for --name, etc.
### 1. Adding a SCM provider
> 1 SCM provider refers to the group of all the repos which can be accessed by a single PAT.

- Github
```bash
prm add provider my-github --type github --host https://github.com
```
- Harness
```bash
prm add provider harness-smp --type harness --host https://smp.harness.com
```
After this you will be prompted to enter your PAT.

<img width="1430" alt="add provider" src="https://github.com/user-attachments/assets/b799040b-0e75-4509-b630-36ea87b748cc">

### 2. Monitoring your PRs
You are ready to start monitoring your PRs.
To list all the PRs ie open, closed, merged
```bash
prm list prs --state all
```
To list the open PRs
```bash
prm list prs --state open
```
<img width="1400" alt="list pr" src="https://github.com/user-attachments/assets/2dbf978a-c2f1-40d9-a43f-ef9a18d3b717">

You can filter the PRs further by provider type (--type) and provider name (--name).
#### Changing the output format
You can change the default format from table to json or yaml. \
json
```bash
prm list prs --output json
```
yaml
```bash
prm list prs --output yaml
```
<img width="1400" alt="yaml-json" src="https://github.com/user-attachments/assets/15bc2704-0747-4315-bdde-69410e996117">

### 3. List your SCM providers
You can check what all SCM providers have been configured.
```bash
prm list providers 
```
<img width="1400" alt="list providers" src="https://github.com/user-attachments/assets/f9ba761f-5f25-41db-abbe-511868ead4c3">

You can filter by name and type.
### 4. Removing an SCM provider
To remove a provider which is no longer needed or is out of date
```bash
prm remove provider my-work-github
```
<img width="1400" alt="remove-provider" src="https://github.com/user-attachments/assets/881c3334-c7f4-49b4-9f97-61b7e045a9d4">

### 5. Refreshing the SCM providers data
> This is applicable only to Harness.
When you add a Harness SCM provider, `prm` fetches user and repo related data which it uses to fetch the PRs. This user and repo data is persisted in a file to reduce unnecessary calls during fetching the PRs. If any org, project or repo has been added or removed for the user, we need to refresh the `prm` config.
```bash
prm refresh providers
```
You can filter by name and type.
### 6. Purging all the SCM providers data saved by prm
If you wish to remove all the data persisted by `prm`
```bash
prm purge
```
You will be prompted for confirmation post running this command.

<img width="1400" alt="purge" src="https://github.com/user-attachments/assets/229bb5d3-119a-411f-a5d1-d371e62c13dc">

Or you can force it
```bash
prm purge --force
```

## Uninstallation
If you want to uninstall prm, you can execute the following
```bash
curl -L https://raw.githubusercontent.com/dhruv1397/prm/main/uninstall.sh | bash
```
<img width="1400" alt="uninstall" src="https://github.com/user-attachments/assets/7e7b50da-c1b1-42c8-925d-d921a1c1de69">


> If `prm` was installed in a system directory like /usr/local/bin, you might need sudo to uninstall it

## Security
It is important to be aware of what data is accessed by any tool to ensure there is no abuse.
To ensure your data is secure, `prm` does not share your data outside your setup.
Moreover, it provides the `purge` command to delete all the data persisted by the app.

## Supported configurations (OS/Arch)
- linux/amd64
- linux/arm64
- darwin/amd64
- darwin/arm64

Check your OS
```bash
uname -s | tr '[:upper:]' '[:lower:]'
```
Check your Arch
```bash
uname -m
```
## References
### Github
Scopes required for PAT: `read:org, repo`

Generate PAT (classic):
https://docs.github.com/en/rest/authentication/authenticating-to-the-rest-api?apiVersion=2022-11-28#authenticating-with-a-personal-access-token \
https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#personal-access-tokens-classic

API:
https://docs.github.com/en/rest?apiVersion=2022-11-28

### Harness
Generate PAT:
https://developer.harness.io/docs/platform/automation/api/add-and-manage-api-keys/#create-personal-api-keys-and-tokens

API:
https://apidocs.harness.io
