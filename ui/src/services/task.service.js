import axios from "axios";

const PIPELINE_API_URL = "/pipeline-api/";

export const TASK_TYPE_QUERY = "query"
export const TASK_TYPE_CONDITION = "condition"
export const TASK_TYPE_AGGREGATOR = "aggregator"
export const TASK_TYPE_TRANSFORMER = "transformer"
export const TASK_TYPE_NOTIFICATION = "notification"
export const TASK_TYPE_API_CALL = "api_call"
export const TASK_TYPE_FOR_EACH = "for_each"
export const TASK_TYPE_MERGE = "merge"
export const TASK_TYPE_FILTER = "filter"
export const TASK_TYPE_EXTERNAL_TRIGGER = "external_trigger"
export const TASK_TYPE_RUN_PIPELINE = "run_pipeline"
export const TASK_TYPE_CODE = "code"
export const TASK_TYPE_PYTHON = "python"
export const TASK_TYPE_GPT = "gpt"
export const TASK_TYPE_PROCESSOR = "processor"

export const TASKS = [
    TASK_TYPE_QUERY,
    TASK_TYPE_CONDITION,
    TASK_TYPE_AGGREGATOR,
    TASK_TYPE_TRANSFORMER,
    TASK_TYPE_NOTIFICATION,
    TASK_TYPE_API_CALL,
    TASK_TYPE_FOR_EACH,
    TASK_TYPE_MERGE,
    TASK_TYPE_FILTER,
    TASK_TYPE_EXTERNAL_TRIGGER,
    TASK_TYPE_RUN_PIPELINE,
    TASK_TYPE_CODE,
];

export const TASK_SEVERITY_LOWEST = "lowest"
export const TASK_SEVERITY_LOW = "low"
export const TASK_SEVERITY_MEDIUM = "medium"
export const TASK_SEVERITY_HIGH = "high"
export const TASK_SEVERITY_CRITICAL = "critical"

export const TASK_PROCESSOR_STRATEGY_INCLUSIVE = "inclusive"
export const TASK_PROCESSOR_STRATEGY_EXCLUSIVE = "exclusive"

// @deprecated. Only notifications have severities, should be retrievable from backend.
export const TASK_SEVERITIES = [
    TASK_SEVERITY_LOWEST,
    TASK_SEVERITY_LOW,
    TASK_SEVERITY_MEDIUM,
    TASK_SEVERITY_HIGH,
    TASK_SEVERITY_CRITICAL,
];

class TaskService {
    getTasksByPipeline = async(uuid) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                PIPELINE_API_URL + 'pipeline/' + uuid + '/tasks',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    searchTasks = async(organizationUuid, string) => {
        var token = localStorage.getItem("token");

        return axios
            .get(
                PIPELINE_API_URL + 'organization/' + organizationUuid + '/tasks/search/' + string,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    deleteTask(uuid, pipelineUuid) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + pipelineUuid + '/task/' + uuid + '/delete',
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    getTask(uuid, pipelineUuid) {
        var token = localStorage.getItem("token")
        return axios
            .get(
                PIPELINE_API_URL + 'pipeline/' + pipelineUuid + '/task/' + uuid,
                {},
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    updateTask(uuid, pipelineUuid, name, severity, type, implementation = null) {
        var token = localStorage.getItem("token");
        var data = { name, severity };

        if (implementation !== null) {
            data[type] = implementation;
        }

        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + pipelineUuid + '/task/' + uuid,
                data,
                { headers: { Authorization: 'Bearer ' + token } }
            );
    }

    addTask(pipelineUuid, name, type) {
        var token = localStorage.getItem("token")
        return axios
            .post(
                PIPELINE_API_URL + 'pipeline/' + pipelineUuid + '/task',
                { name, type },
                { headers: { Authorization: 'Bearer ' + token } }
            )
            .then((response) => {
                return response.data;
            });
    }
}

export default new TaskService();
