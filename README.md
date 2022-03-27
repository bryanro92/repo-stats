 # repo-stats

 * This project looks at a repository, and given a timeframe of analysis, will review the total amount of activity on pull requests within that timeframe.

## Setup

* Create a folder called secrets

```
mkdir secrets
```

```
cat > secrets/env <<'EOF'
GH_PAT=<your token here>
EOF
```

* `GH_PAT` is your personal access token from GitHub. You can use this anonymously, but stricter rate limits apply.

## Build

```
make build
```

## Run

```
make run
```

## Parameters

1. `owner` is the organization or user for the repository.

1. `repository` is the repository we will be examininig.

1. `days` the number of days to look back at.

## Example

* Make run

```
make run
```

* Go Run

```
go run . azure aro-rp 7
```