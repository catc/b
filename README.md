# b
Git branch tools + shortcuts. This tool was created for own personal workflow - your mileage may vary, use at own risk.

<img src="./assets/screenshot.png" />

## Usage

#### `b`
View available branches to checkout
<br/>

#### `b -l`
View local branches in repo
<br/>

#### `b prune`
Select branches to delete
<br/>

#### `b prune 20`
Select branches to delete that are older than n days
<br/>

#### `b clone`
Clone branch - useful if about to make potentially destructive changes and don't wanna go through reflog to fix it
<br/>

#### `b add`
Walk through changed + untracked files to stage, basically a glorified git diff with the option to add, skip or checkout unstaged files
<br/>


## Dev
- install deps via `go mod download`
- build and run `go build && ./b`
