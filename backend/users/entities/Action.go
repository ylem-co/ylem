package entities

const ACTION_CREATE = "create"
const ACTION_READ = "read"
const ACTION_READ_LIST = "read_list"
const ACTION_UPDATE = "update"
const ACTION_DELETE = "delete"
const ACTION_RUN = "run"

func IsActionValid(action string) bool {
    switch action {
    case
        ACTION_CREATE,
        ACTION_READ,
        ACTION_READ_LIST,
        ACTION_UPDATE,
        ACTION_RUN,
        ACTION_DELETE:
        return true
    }
    return false
}
